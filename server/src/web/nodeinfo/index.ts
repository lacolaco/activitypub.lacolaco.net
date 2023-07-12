import { Handler, Hono } from 'hono';
import * as pkg from '../../../../package.json';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/nodeinfo', handleNodeinfo);
  app.get('/nodeinfo/2.1', handleNodeinfo21);
};

const handleNodeinfo: Handler = async (c) => {
  const origin = c.get('origin');

  const res = c.json({
    links: [
      {
        rel: 'http://nodeinfo.diaspora.software/ns/schema/2.1',
        href: `${origin}/nodeinfo/2.1`,
      },
    ],
  });
  return res;
};

const handleNodeinfo21: Handler = async (c) => {
  const res = c.json({
    version: '2.1',
    openRegistrations: false,
    protocols: ['activitypub'],
    software: {
      name: 'Hono',
      version: pkg.dependencies.hono,
    },
    usage: {
      users: {
        total: 1,
      },
    },
    services: {
      inbound: [],
      outbound: [],
    },
    metadata: {
      nodeName: 'らこらこインターネット',
      nodeDescription: 'らこらこインターネット',
      disableRegistration: true,
      themeColor: '#77b58c',
      maintainer: {
        name: 'lacolaco',
        email: 'https://github.com/lacolaco#where-you-can-contact-me',
      },
    },
  });
  return res;
};
