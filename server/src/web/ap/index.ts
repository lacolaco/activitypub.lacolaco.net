import * as ap from '@app/activitypub';
import { User } from '@app/domain/user';
import { UsersRepository } from '@app/repository/users';
import { runInSpan } from '@app/tracing';
import { acceptFollowRequest, deleteFollower, getUserFollowers } from '@app/usecase/followers';
import { Handler, Hono, MiddlewareHandler } from 'hono';
import { assertContentTypeHeader } from '../../middleware/asserts';
import { AppContext } from '../context';

type UserRouteContext = AppContext & {
  Variables: {
    user: User;
  };
};

function setUserMiddleware(): MiddlewareHandler<UserRouteContext> {
  return async (c, next) => {
    const userRepo = new UsersRepository();
    const id = c.req.param('id');
    const user = await userRepo.findByID(id);
    if (user == null) {
      // if an user has the username, tell client to redirect permanently
      const origin = c.get('origin');
      const hostname = new URL(origin).hostname;
      const u = await userRepo.findByUsername(hostname, id);
      if (u != null) {
        const redirectTo = new URL(origin);
        redirectTo.pathname = c.req.path.replace(id, u.id);
        console.log('redirecting to', redirectTo.toString());
        c.res.headers.set('Location', redirectTo.toString());
        return c.json({ error: 'Moved Permanently' }, 301);
      }
      return c.json({ error: 'Not Found' }, 404);
    }
    c.set('user', user);
    await next();
  };
}

const setActivityJSONContentType = (): MiddlewareHandler => async (c, next) => {
  await next();
  c.res.headers.set('Content-Type', 'application/activity+json');
};

export default (app: Hono<AppContext>) => {
  const apRoutes = new Hono<AppContext>();
  // middlewares
  apRoutes.get('*', setActivityJSONContentType());
  apRoutes.post('*', assertContentTypeHeader(['application/activity+json']));
  // routes
  apRoutes.get('/inbox', handleGetSharedInbox);
  apRoutes.post('/inbox', handlePostSharedInbox);

  const userRoutes = new Hono<UserRouteContext>();
  userRoutes.use('*', setUserMiddleware());
  userRoutes.get('/', handleGetPerson);
  userRoutes.get('/inbox', handleGetInbox);
  userRoutes.post('/inbox', handlePostInbox);
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
    const person = ap.buildPerson(origin, user, config.publicKeyPem);
    const res = c.json(person);
    return res;
  });
};

const handleGetInbox: Handler<UserRouteContext> = async (c) => {
  return c.json({}, 405);
};

const handlePostInbox: Handler<UserRouteContext> = async (c) => {
  return runInSpan('ap.handlePostInbox', async (span) => {
    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      console.error(e);
      return c.json({ error: 'Unauthorized' }, 401);
    }
    const config = c.get('config');
    const origin = c.get('origin');
    const user = c.get('user');

    const payload = await c.req.json();
    const parsed = ap.AnyActivity.safeParse(payload);
    if (!parsed.success) {
      console.error(JSON.stringify(parsed.error));
      return c.json({ error: 'Bad Request' }, 400);
    }
    const activity = parsed.data;
    console.debug(JSON.stringify(activity));

    if (ap.isFollowActivity(activity)) {
      try {
        await acceptFollowRequest(config, origin, user, activity);
        return c.json({ ok: true });
      } catch (e) {
        console.error(e);
        return c.json({ error: 'Internal Server Error' }, 500);
      }
    } else if (ap.isUndoActivity(activity)) {
      const object = activity.object;
      // unfollow
      if (ap.isFollowActivity(object)) {
        try {
          await deleteFollower(config, origin, user, activity);
          return c.json({ ok: true });
        } catch (e) {
          console.error(e);
          return c.json({ error: 'Internal Server Error' }, 500);
        }
      }
    }

    console.log('unsupported activity');
    return c.json({});
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
    return c.json(
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
    return c.json(ap.buildOrderedCollection(person.following, []));
  });
};

const handleGetSharedInbox: Handler<AppContext> = async (c) => {
  return c.json({}, 405);
};

const handlePostSharedInbox: Handler<AppContext> = async (c) => {
  return runInSpan('ap.handlePostSharedInbox', async (span) => {
    const payload = await c.req.json();
    console.debug(JSON.stringify(payload));

    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      console.error(e);
      return c.json({ error: 'Unauthorized' }, 401);
    }

    const parsed = ap.AnyActivity.safeParse(payload);
    if (!parsed.success) {
      console.error(JSON.stringify(parsed.error));
      return c.json({ error: 'Bad Request' }, 400);
    }
    const activity = parsed.data;

    span.setAttributes({
      'activity.type': activity.type,
      'activity.actor': ap.getURI(activity.actor)?.toString(),
    });

    return c.json({ error: 'Not Found' }, 404);
  });
};
