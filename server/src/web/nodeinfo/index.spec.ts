import { Hono } from 'hono';
import { describe, expect, test } from 'vitest';
import { AppContext, withOrigin } from '../context';
import useWebfinger from './index';

describe('nodeinfo', () => {
  const app = new Hono<AppContext>();
  app.use('*', withOrigin());
  useWebfinger(app);

  test('entrypoint response is valid', async () => {
    const req = new Request('http://localhost/.well-known/nodeinfo', {
      method: 'GET',
      headers: { host: 'localhost' },
    });
    const res = await app.request(req);
    expect(res.headers.get('Content-Type')).toContain('application/json');
    const body = await res.json();

    expect(body).toMatchInlineSnapshot(`
      {
        "links": [
          {
            "href": "https://localhost/nodeinfo/2.1",
            "rel": "http://nodeinfo.diaspora.software/ns/schema/2.1",
          },
        ],
      }
    `);
  });

  test('nodeinfo response is valid', async () => {
    const req = new Request('http://localhost/nodeinfo/2.1', {
      method: 'GET',
      headers: { host: 'localhost' },
    });
    const res = await app.request(req);
    expect(res.headers.get('Content-Type')).toContain('application/json');
    const body = await res.json();

    expect(body).toMatchInlineSnapshot(`
      {
        "metadata": {
          "disableRegistration": true,
          "maintainer": {
            "email": "https://github.com/lacolaco#where-you-can-contact-me",
            "name": "lacolaco",
          },
          "nodeDescription": "らこらこインターネット",
          "nodeName": "らこらこインターネット",
          "themeColor": "#77b58c",
        },
        "openRegistrations": false,
        "protocols": [
          "activitypub",
        ],
        "services": {
          "inbound": [],
          "outbound": [],
        },
        "software": {
          "name": "Hono",
          "version": "^3.3.0",
        },
        "usage": {
          "users": {
            "total": 1,
          },
        },
        "version": "2.1",
      }
    `);
  });
});
