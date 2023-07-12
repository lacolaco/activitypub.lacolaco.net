import * as ap from '@app/activitypub';
import { UsersRepository } from '@app/repository/users';
import { Handler, Hono, MiddlewareHandler } from 'hono';
import { acceptFollowRequest, deleteFollower, getUserFollowers } from 'server/src/usecase/relationship';
import { assertContentTypeHeader } from '../../middleware/asserts';
import { AppContext } from '../context';

const setActivityJSONContentType = (): MiddlewareHandler => async (c, next) => {
  await next();
  c.res.headers.set('Content-Type', 'application/activity+json');
};

export default (app: Hono<AppContext>) => {
  const apRoutes = new Hono();
  // middlewares
  apRoutes.get('*', setActivityJSONContentType());
  apRoutes.post('*', assertContentTypeHeader(['application/activity+json']));
  // routes
  apRoutes.post('/inbox', handlePostSharedInbox);

  const userRoutes = new Hono();
  userRoutes.get('/', handleGetPerson);
  userRoutes.post('/inbox', handlePostInbox);
  userRoutes.get('/outbox', handleGetOutbox);
  userRoutes.get('/followers', handleGetFollowers);
  userRoutes.get('/following', handleGetFollowing);

  apRoutes.route('/users/:id', userRoutes);
  app.route('/', apRoutes);
};

const handleGetPerson: Handler<AppContext> = async (c) => {
  const { origin } = new URL(c.req.url);

  const userRepo = new UsersRepository();
  const id = c.req.param('id');

  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }
  const person = ap.setPublicKey(ap.buildPerson(origin, user), c.get('rsaKeyPair').publicKey);
  const res = c.json(person);
  return res;
};

const handlePostInbox: Handler<AppContext> = async (c) => {
  try {
    await ap.verifySignature(c.req);
  } catch (e) {
    c.status(400);
    return c.text('Bad Request');
  }

  const userRepo = new UsersRepository();
  const id = c.req.param('id');
  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }

  const activity = await c.req.json<ap.Activity>();
  if (ap.isFollowActivity(activity)) {
    try {
      const config = c.get('Config');
      await acceptFollowRequest(user, activity, config);
      return c.json({ ok: true });
    } catch (e) {
      console.error(e);
      c.status(500);
      return c.text('Internal Server Error');
    }
  } else if (ap.isUndoActivity(activity)) {
    const object = activity.object;
    // unfollow
    if (ap.isFollowActivity(object)) {
      try {
        const config = c.get('Config');
        await deleteFollower(user, activity, config);
        return c.json({ ok: true });
      } catch (e) {
        console.error(e);
        c.status(500);
        return c.text('Internal Server Error');
      }
    }
  }

  console.log('unsupported activity');
  return c.text('ok');
};

const handleGetOutbox: Handler = async (c) => {
  const { origin } = new URL(c.req.url);
  const userRepo = new UsersRepository();
  const id = c.req.param('id');

  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.outbox, []));
  return res;
};

const handleGetFollowers: Handler = async (c) => {
  const { origin } = new URL(c.req.url);

  const userRepo = new UsersRepository();
  const id = c.req.param('id');
  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }

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

const handleGetFollowing: Handler = async (c) => {
  const { origin } = new URL(c.req.url);
  const userRepo = new UsersRepository();
  const id = c.req.param('id');

  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.following, []));
  return res;
};

const handlePostSharedInbox: Handler = async (c) => {
  return c.text('ok');
};
