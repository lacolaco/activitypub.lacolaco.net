import { UsersRepository } from '@app/repository/users';
import { JRDObject } from '@app/webfinger';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/webfinger', handleWebfinger);
};

const handleWebfinger: Handler = async (c) => {
  const origin = c.get('origin');

  // resource format: acct:username@domain
  const resource = c.req.query('resource');
  if (resource == null) {
    return c.json({ error: 'Bad Request' }, 400);
  }
  const [, username] = resource.split('@')[0].split(':');

  const userRepo = new UsersRepository();
  const user = await userRepo.findByUsername(username);
  if (user == null) {
    return c.json({ error: 'Not Found' }, 404);
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
