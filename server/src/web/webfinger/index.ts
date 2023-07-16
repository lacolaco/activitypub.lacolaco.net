import { UsersRepository } from '@app/repository/users';
import { JRDObject } from '@app/webfinger';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/webfinger', handleWebfinger);
};

const resourceRegexp = /^acct\:(.+)@(.+)$/;

const handleWebfinger: Handler = async (c) => {
  const origin = c.get('origin');

  // resource format: acct:username@domain
  const resource = c.req.query('resource');
  if (resource == null) {
    return c.json({ error: 'Bad Request' }, 400);
  }
  const [, username, hostname] = resource.match(resourceRegexp) ?? [];
  if (username == null || hostname == null) {
    return c.json({ error: 'Bad Request' }, 400);
  }
  const userRepo = new UsersRepository();
  const user = await userRepo.findByUsername(hostname, username);
  if (user == null) {
    return c.json({ error: 'Not Found' }, 404);
  }

  return c.json<JRDObject>(
    {
      subject: resource,
      links: [
        {
          rel: 'self',
          type: 'application/activity+json',
          href: `${origin}/users/${user.id}`,
        },
      ],
    },
    { headers: { 'Content-Type': 'application/jrd+json' } },
  );
};
