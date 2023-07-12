import { Hono } from 'hono';
import { assert, beforeEach, describe, test } from 'vitest';

import createApplication from './app';
import { AppContext } from './web/context';

describe('endpoints for activitypub compatibility', () => {
  let app: Hono<AppContext>;

  beforeEach(async () => {
    app = await createApplication();
  });

  test('webfinger is supported', async () => {
    const req = new Request('http://localhost/.well-known/webfinger?resource=acct:alice@localhost', {
      method: 'GET',
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert.equal(res.headers.get('Content-Type'), 'application/jrd+json');
  });

  test('nodeinfo is supported', async () => {
    const req = new Request('http://localhost/.well-known/nodeinfo', {
      method: 'GET',
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/json'));
  });

  test('host-meta (xml) is supported', async () => {
    const req = new Request('http://localhost/.well-known/host-meta', {
      method: 'GET',
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/xrd+xml'));
  });

  test('host-meta (json) is supported', async () => {
    const req = new Request('http://localhost/.well-known/host-meta.json', {
      method: 'GET',
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/json'));
  });
});
