import { Handler, Hono } from 'hono';
import * as pkg from '../../../../package.json';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/nodeinfo', handleNodeinfo);
  app.get('/nodeinfo/2.1', handleNodeinfo21);
};

const handleNodeinfo: Handler = async (c) => {
  const origin = c.get('origin');
  return c.json(
    {
      links: [
        {
          rel: 'http://nodeinfo.diaspora.software/ns/schema/2.1',
          href: `${origin}/nodeinfo/2.1`,
        },
      ],
    },
    {
      headers: {
        'Cache-Control': `max-age=${60 * 60}, public`,
      },
    },
  );
};

const handleNodeinfo21: Handler = async (c) => {
  return c.json(
    {
      version: '2.1',
      openRegistrations: false,
      protocols: ['activitypub'],
      software: {
        name: 'Hono',
        version: pkg.dependencies.hono,
      },
      usage: {
        users: { total: 1 },
      },
      services: { inbound: [], outbound: [] },
      metadata: {
        nodeName: 'lacolaco.social',
        nodeDescription: 'らこらこインターネット',
        disableRegistration: true,
        themeColor: '#77b58c',
        maintainer: {
          name: 'lacolaco',
          email: 'https://github.com/lacolaco#where-you-can-contact-me',
        },
      },
    },
    {
      headers: {
        'Cache-Control': `max-age=${60 * 60}, public`,
      },
    },
  );
};
