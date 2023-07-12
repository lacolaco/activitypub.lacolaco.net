import { Handler, Hono } from 'hono';
import { JRDObject } from './types';
import { AppContext } from '../context';
import { UsersRepository } from '@app/repository/users';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/webfinger', handleWebfinger);
};

const handleWebfinger: Handler = async (c) => {
  // resource format: acct:username@domain
  const resource = c.req.query('resource');
  if (resource == null) {
    c.status(400);
    return c.json({ error: 'Bad Request' });
  }
  const { origin } = new URL(c.req.url);
  const [, username] = resource.split('@')[0].split(':');

  const userRepo = new UsersRepository();
  const user = await userRepo.findByUsername(username);
  if (user == null) {
    c.status(404);
    return c.json({ error: 'Not Found' });
  }

  const res = c.json<JRDObject>({
    subject: resource,
    links: [
      {
        rel: 'self',
        type: 'application/activity+json',
        href: `${origin}/users/${user.id}`,
      },
    ],
  });
  res.headers.set('Content-Type', 'application/jrd+json');
  return res;
};
