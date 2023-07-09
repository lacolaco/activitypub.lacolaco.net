import * as ap from '@app/activitypub';
import { UsersRepository } from '@app/repository/users';
import { Handler, Hono, MiddlewareHandler } from 'hono';
import { Env } from '../../env';
import { assertContentTypeHeader } from '../../middleware/asserts';

const setContentType = (): MiddlewareHandler => async (c, next) => {
	await next();
	c.res.headers.set('Content-Type', 'application/activity+json');
};

export default (app: Hono) => {
	const apRoutes = new Hono<Env>();
	// middlewares
	apRoutes.get('*', setContentType());
	apRoutes.post('*', assertContentTypeHeader(['application/activity+json']));
	// routes
	apRoutes.post('/inbox', handlePostSharedInbox);

	const userRoutes = new Hono<Env>();
	userRoutes.get('/', handleGetPerson);
	userRoutes.post('/inbox', handlePostInbox);
	userRoutes.get('/outbox', handleGetOutbox);
	userRoutes.get('/followers', handleGetFollowers);
	userRoutes.get('/following', handleGetFollowing);

	apRoutes.route('/users/:username', userRoutes);
	app.route('/ap', apRoutes);
};

const handleGetPerson: Handler<Env> = async (c) => {
	const { origin } = new URL(c.req.url);
	const userRepo = new UsersRepository(c.env.DB);
	const username = c.req.param('username');

	const user = await userRepo.findByUsername(username);
	if (user == null) {
		c.status(404);
		return c.text('Not Found');
	}
	const res = c.json(ap.buildPerson(origin, user));
	return res;
};

const handlePostInbox: Handler<Env> = async (c) => {
	return c.text('ok');
};

const handleGetOutbox: Handler<Env> = async (c) => {
	const { origin } = new URL(c.req.url);
	const userRepo = new UsersRepository(c.env.DB);
	const username = c.req.param('username');

	const user = await userRepo.findByUsername(username);
	if (user == null) {
		c.status(404);
		return c.text('Not Found');
	}
	const person = ap.buildPerson(origin, user);
	const res = c.json(ap.buildOrderedCollection(person.outbox, []));
	return res;
};

const handleGetFollowers: Handler<Env> = async (c) => {
	const { origin } = new URL(c.req.url);
	const userRepo = new UsersRepository(c.env.DB);
	const username = c.req.param('username');

	const user = await userRepo.findByUsername(username);
	if (user == null) {
		c.status(404);
		return c.text('Not Found');
	}
	const person = ap.buildPerson(origin, user);
	const res = c.json(ap.buildOrderedCollection(person.followers, []));
	return res;
};

const handleGetFollowing: Handler<Env> = async (c) => {
	const { origin } = new URL(c.req.url);
	const userRepo = new UsersRepository(c.env.DB);
	const username = c.req.param('username');

	const user = await userRepo.findByUsername(username);
	if (user == null) {
		c.status(404);
		return c.text('Not Found');
	}
	const person = ap.buildPerson(origin, user);
	const res = c.json(ap.buildOrderedCollection(person.following, []));
	return res;
};

const handlePostSharedInbox: Handler<Env> = async (c) => {
	return c.text('ok');
};
