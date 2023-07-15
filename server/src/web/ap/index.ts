import * as ap from '@app/activitypub';
import { User } from '@app/domain/user';
import { UsersRepository } from '@app/repository/users';
import { getTracer } from '@app/tracing';
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
      const u = await userRepo.findByUsername(id);
      if (u != null) {
        c.status(301);
        const redirectTo = new URL(c.get('origin'));
        redirectTo.pathname = c.req.path.replace(id, u.id);
        console.log('redirecting to', redirectTo.toString());
        c.res.headers.set('Location', redirectTo.toString());
        return c.json({ error: 'Moved Permanently' });
      }

      c.status(404);
      return c.json({ error: 'Not Found' });
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
  return getTracer().startActiveSpan('ap.handleGetPerson', async (span) => {
    const config = c.get('config');
    const origin = c.get('origin');
    const user = c.get('user');
    const person = ap.withPublicKey(ap.buildPerson(origin, user), config.publicKeyPem);
    const res = c.json(person);
    return res;
  });
};

const handleGetInbox: Handler<UserRouteContext> = async (c) => {
  const origin = c.get('origin');
  const user = c.get('user');
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.inbox, []));
  return res;
};

const handlePostInbox: Handler<UserRouteContext> = async (c) => {
  return getTracer().startActiveSpan('ap.handlePostInbox', async (span) => {
    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      console.error(e);
      c.status(400);
      return c.json({ error: 'Bad Request' });
    }
    const config = c.get('config');
    const origin = c.get('origin');
    const user = c.get('user');

    const activity = await c.req.json<ap.Activity>();
    console.debug(JSON.stringify(activity));

    if (ap.isFollowActivity(activity)) {
      try {
        await acceptFollowRequest(config, origin, user, activity);
        return c.json({ ok: true });
      } catch (e) {
        console.error(e);
        c.status(500);
        return c.json({ error: 'Internal Server Error' });
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
          c.status(500);
          return c.json({ error: 'Internal Server Error' });
        }
      }
    }

    console.log('unsupported activity');
    return c.json({});
  });
};

const handleGetOutbox: Handler<UserRouteContext> = async (c) => {
  const origin = c.get('origin');
  const user = c.get('user');
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.outbox, []));
  return res;
};

const handleGetFollowers: Handler<UserRouteContext> = async (c) => {
  const origin = c.get('origin');
  const user = c.get('user');

  const followers = await getUserFollowers(user);

  const person = ap.buildPerson(origin, user);
  const res = c.json(
    ap.buildOrderedCollection(
      person.followers,
      followers.map((f) => new URL(f.id)),
    ),
  );
  return res;
};

const handleGetFollowing: Handler<UserRouteContext> = async (c) => {
  const origin = c.get('origin');
  const user = c.get('user');
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.following, []));
  return res;
};

const handleGetSharedInbox: Handler<AppContext> = async (c) => {
  c.status(404);
  return c.json({ error: 'Not Found' });
};

const handlePostSharedInbox: Handler<AppContext> = async (c) => {
  return getTracer().startActiveSpan('ap.handlePostSharedInbox', async (span) => {
    try {
      await ap.verifySignature(c.req);
    } catch (e) {
      console.error(e);
      c.status(400);
      return c.json({ error: 'Bad Request' });
    }

    const activity = await c.req.json<ap.Activity>();
    console.debug(JSON.stringify(activity));

    span.setAttributes({
      'activity.type': activity.type,
      'activity.actor': ap.getID(activity.actor)?.toString(),
    });

    c.status(404);
    return c.json({ error: 'Not Found' });
  });
};
