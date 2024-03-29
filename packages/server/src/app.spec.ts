import { Hono } from 'hono';
import { assert, beforeEach, describe, expect, test } from 'vitest';

import createApplication from './app';
import { getConfigWithEnv } from './domain/config';
import { AppContext } from './web/context';

describe('endpoints for activitypub compatibility', () => {
  let app: Hono<AppContext>;

  beforeEach(async () => {
    const config = getConfigWithEnv();
    app = await createApplication(config);
  });

  test('all routes are registered', async () => {
    expect(app.routes.map((r) => `${r.method} ${r.path} ${r.handler.name || '[inline]'}`)).toMatchInlineSnapshot(`
      [
        "ALL * [inline]",
        "ALL * [inline]",
        "ALL * [inline]",
        "ALL * [inline]",
        "ALL * [inline]",
        "GET /.well-known/nodeinfo handleNodeinfo",
        "GET /nodeinfo/2.1 handleNodeinfo21",
        "GET /.well-known/host-meta handleHostMetaXML",
        "GET /.well-known/host-meta.json handleHostMetaJSON",
        "GET /.well-known/webfinger handleWebfinger",
        "GET /inbox handleGetSharedInbox",
        "POST /inbox [inline]",
        "POST /inbox handlePostSharedInbox",
        "ALL /users/:id/* [inline]",
        "GET /users/:id handleGetPerson",
        "GET /users/:id/inbox handleGetInbox",
        "POST /users/:id/inbox [inline]",
        "POST /users/:id/inbox handlePostInbox",
        "GET /users/:id/outbox handleGetOutbox",
        "GET /users/:id/followers handleGetFollowers",
        "GET /users/:id/following handleGetFollowing",
        "GET /admin/users [inline]",
        "POST /admin/users [inline]",
        "GET /admin/users/:hostname/:username [inline]",
        "POST /admin/users/:hostname/:username/notes [inline]",
        "GET /admin/search/person/:resource handleSearchPerson",
      ]
    `);
  });

  // TODO: UsersRepository をモックできるようにする
  test.skip('webfinger is supported', async () => {
    const req = new Request('http://localhost/.well-known/webfinger?resource=acct:alice@localhost', {
      method: 'GET',
      headers: { host: 'localhost:80' },
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert.equal(res.headers.get('Content-Type'), 'application/jrd+json');
  });

  test('nodeinfo is supported', async () => {
    const req = new Request('http://localhost/.well-known/nodeinfo', {
      method: 'GET',
      headers: { host: 'localhost:80' },
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/json'));
  });

  test('host-meta (xml) is supported', async () => {
    const req = new Request('http://localhost/.well-known/host-meta', {
      method: 'GET',
      headers: { host: 'localhost:80' },
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/xrd+xml'));
  });

  test('host-meta (json) is supported', async () => {
    const req = new Request('http://localhost/.well-known/host-meta.json', {
      method: 'GET',
      headers: { host: 'localhost:80' },
    });
    const res = await app.request(req);
    assert.equal(res.status, 200);
    assert(res.headers.get('Content-Type')?.includes('application/json'));
  });
});
