import { Handler, Hono } from 'hono';
import { JRDObject } from './types';

export default (app: Hono) => {
  app.get('/.well-known/webfinger', handleWebfinger);
};

const handleWebfinger: Handler = async (c) => {
  // resource format: acct:username@domain
  const resource = c.req.query('resource');
  if (resource == null) {
    c.status(400);
    return c.text('Bad Request');
  }
  const { origin } = new URL(c.req.url);
  const username = resource.split('@')[0].split(':')[1];

  const res = c.json<JRDObject>({
    subject: resource,
    links: [
      {
        rel: 'self',
        type: 'application/activity+json',
        href: `${origin}/ap/users/${username}`,
      },
    ],
  });
  res.headers.set('Content-Type', 'application/jrd+json');
  return res;
};
