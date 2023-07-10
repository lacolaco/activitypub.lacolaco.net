import * as ap from '@app/activitypub';
import { UsersRepository } from '@app/repository/users';
import { Handler, Hono, MiddlewareHandler } from 'hono';
import { assertContentTypeHeader } from '../../middleware/asserts';

const setContentType = (): MiddlewareHandler => async (c, next) => {
  await next();
  c.res.headers.set('Content-Type', 'application/activity+json');
};

export default (app: Hono) => {
  const apRoutes = new Hono();
  // middlewares
  apRoutes.get('*', setContentType());
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

const handleGetPerson: Handler = async (c) => {
  const { origin } = new URL(c.req.url);
  const userRepo = new UsersRepository();
  const id = c.req.param('id');

  const user = await userRepo.findByID(id);
  if (user == null) {
    c.status(404);
    return c.text('Not Found');
  }
  const res = c.json(ap.buildPerson(origin, user));
  return res;
};

const handlePostInbox: Handler = async (c) => {
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
  const person = ap.buildPerson(origin, user);
  const res = c.json(ap.buildOrderedCollection(person.followers, []));
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
