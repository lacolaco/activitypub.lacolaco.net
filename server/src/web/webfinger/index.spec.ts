import { Hono } from 'hono';
import { describe, test, assert, expect } from 'vitest';
import useWebfinger from './index';
import { JRDObject } from './types';

describe('webfinger', () => {
	const app = new Hono();
	useWebfinger(app);

	test('response is JRD content', async () => {
		const req = new Request('http://localhost/.well-known/webfinger?resource=acct:alice@localhost', {
			method: 'GET',
		});
		const res = await app.request(req);
		assert.equal(res.headers.get('Content-Type'), 'application/jrd+json');
		const body = await res.json<JRDObject>();
		assert.equal(body.subject, 'acct:alice@localhost');
		assert.equal(body.links[0].rel, 'self');
		assert.equal(body.links[0].type, 'application/activity+json');
		assert.equal(body.links[0].href, 'http://localhost/ap/users/alice');

		expect(body).toMatchInlineSnapshot(`
			{
			  "links": [
			    {
			      "href": "http://localhost/ap/users/alice",
			      "rel": "self",
			      "type": "application/activity+json",
			    },
			  ],
			  "subject": "acct:alice@localhost",
			}
		`);
	});
});
