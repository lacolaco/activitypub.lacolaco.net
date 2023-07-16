import * as ap from '@app/activitypub';
import { User } from '@app/domain/user';
import { UsersRepository } from '@app/repository/users';
import { runInSpan } from '@app/tracing';
import { acceptFollowRequest, deleteFollower, getUserFollowers } from '@app/usecase/followers';
import { Context, Handler, Hono, MiddlewareHandler } from 'hono';
import { assertContentTypeHeader } from '../assertion';
import { AppContext } from '../context';

type UserRouteContext = AppContext & {
  Variables: {
    user: User;
  };
};

function setUserMiddleware(): MiddlewareHandler<UserRouteContext> {
  return async (c, next) => {
    await runInSpan('ap.setUserMiddleware', async (span) => {
      const logger = c.get('logger');
      const userRepo = new UsersRepository();
      const id = c.req.param('id');
      span.setAttribute('request.id', id);
      const user = await userRepo.findByID(id);
      if (user == null) {
        // if an user has the username, tell client to redirect permanently
        const origin = c.get('origin');
        const hostname = new URL(origin).hostname;
        const u = await userRepo.findByUsername(hostname, id);
        if (u != null) {
          const redirectTo = new URL(origin);
          redirectTo.pathname = c.req.path.replace(id, u.id);
          logger.info('redirecting to', redirectTo.toString());
          c.res.headers.set('Location', redirectTo.toString());
          return c.json({ error: 'Moved Permanently' }, 301);
        }
        return c.json({ error: 'Not Found' }, 404);
      }
      c.set('user', user);
      await next();
    });
  };
}

function activityjson<T>(c: Context, data: T, status = 200) {
  return c.json(data, { status, headers: { 'Content-Type': 'application/activity+json' } });
}

export default (app: Hono<AppContext>) => {
  const assertActivityJSONContentType = assertContentTypeHeader(['application/activity+json']);

  const apRoutes = new Hono<AppContext>();
  // routes
  apRoutes.get('/inbox', handleGetSharedInbox);
  apRoutes.post('/inbox', assertActivityJSONContentType, handlePostSharedInbox);

  const userRoutes = new Hono<UserRouteContext>();
  // https://github.com/honojs/hono/issues/1240
  userRoutes.get('/', setUserMiddleware(), handleGetPerson);
  userRoutes.use('*', setUserMiddleware());
  userRoutes.get('/inbox', handleGetInbox);
  userRoutes.post('/inbox', assertActivityJSONContentType, handlePostInbox);
  userRoutes.get('/outbox', handleGetOutbox);
  userRoutes.get('/followers', handleGetFollowers);
  userRoutes.get('/following', handleGetFollowing);
  apRoutes.route('/users/:id', userRoutes);
  app.route('/', apRoutes);
};

const handleGetPerson: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handleGetPerson', async (span) => {
    const config = c.get('config');
    const origin = c.get('origin');
    const user = c.get('user');
    c.get('logger').debug({ user });
    const person = ap.buildPerson(origin, user, config.publicKeyPem);
    return activityjson(c, person);
  });
};

const handleGetInbox: Handler<UserRouteContext> = async (c) => {
  return c.json({}, 405);
};

const handlePostInbox: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handlePostInbox', async (span) => {
    const logger = c.get('logger');
    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      logger.error(e);
      return c.json({ error: 'Unauthorized' }, 401);
    }
    const config = c.get('config');
    const origin = c.get('origin');
    const user = c.get('user');

    const payload = await c.req.json();
    const parsed = ap.AnyActivity.safeParse(payload);
    if (!parsed.success) {
      logger.error(parsed.error);
      return c.json({ error: 'Bad Request' }, 400);
    }
    const activity = parsed.data;
    logger.debug({ activity });

    if (ap.isFollowActivity(activity)) {
      try {
        await acceptFollowRequest(config, origin, user, activity);
        return activityjson(c, { ok: true });
      } catch (e) {
        logger.error(e);
        return c.json({ error: 'Internal Server Error' }, 500);
      }
    } else if (ap.isUndoActivity(activity)) {
      const object = activity.object;
      // unfollow
      if (ap.isFollowActivity(object)) {
        try {
          await deleteFollower(config, origin, user, activity);
          return activityjson(c, { ok: true });
        } catch (e) {
          logger.error(e);
          return c.json({ error: 'Internal Server Error' }, 500);
        }
      }
    }

    logger.info('unsupported activity');
    return activityjson(c, {});
  });
};

const handleGetOutbox: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handleGetOutbox', async (span) => {
    const origin = c.get('origin');
    const user = c.get('user');
    const person = ap.buildPerson(origin, user);
    const res = c.json(ap.buildOrderedCollection(person.outbox, []));
    return res;
  });
};

const handleGetFollowers: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handleGetFollowers', async (span) => {
    const origin = c.get('origin');
    const user = c.get('user');
    const person = ap.buildPerson(origin, user);
    if (!person.followers) {
      return c.json({ error: 'Not Found' }, 404);
    }

    const followers = await getUserFollowers(user);
    return activityjson(
      c,
      ap.buildOrderedCollection(
        person.followers,
        followers.map((f) => f.id),
      ),
    );
  });
};

const handleGetFollowing: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handleGetFollowing', async (span) => {
    const origin = c.get('origin');
    const user = c.get('user');
    const person = ap.buildPerson(origin, user);
    if (!person.following) {
      return c.json({ error: 'Not Found' }, 404);
    }
    return activityjson(c, ap.buildOrderedCollection(person.following, []));
  });
};

const handleGetSharedInbox: Handler<AppContext> = async (c) => {
  return c.json({}, 405);
};

const handlePostSharedInbox: Handler<AppContext> = async (c) => {
  return runInSpan('ap.handlePostSharedInbox', async (span) => {
    const logger = c.get('logger');
    const payload = await c.req.json();
    logger.debug(JSON.stringify(payload));
    logger.debug({ payload });

    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      logger.error(e);
      return c.json({ error: 'Unauthorized' }, 401);
    }

    const parsed = ap.AnyActivity.safeParse(payload);
    if (!parsed.success) {
      logger.error(JSON.stringify(parsed.error));
      return c.json({ error: 'Bad Request' }, 400);
    }
    const activity = parsed.data;
    logger.debug({ activity });

    span.setAttributes({
      'activity.type': activity.type,
      'activity.actor': ap.getURI(activity.actor)?.toString(),
    });

    return c.json({ error: 'Not Found' }, 404);
  });
};
