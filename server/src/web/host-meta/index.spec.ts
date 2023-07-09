import { Hono } from 'hono';
import { describe, expect, test } from 'vitest';
import useHostmeta from './index';

describe('hostmeta', () => {
	const app = new Hono();
	useHostmeta(app);

	test('response has a lrdd link', async () => {
		const req = new Request('http://localhost/.well-known/host-meta', {
			method: 'GET',
		});
		const res = await app.request(req);
		const body = await res.text();
		expect(body).toMatchInlineSnapshot(`
			"<?xml version=\\"1.0\\" encoding=\\"UTF-8\\"?>
			<XRD xmlns=\\"http://docs.oasis-open.org/ns/xri/xrd-1.0\\">
			    <Link rel=\\"lrdd\\" template=\\"http://localhost/.well-known/webfinger?resource={uri}\\"/>
			</XRD>"
		`);
	});

	test('response has a lrdd link (json)', async () => {
		const req = new Request('http://localhost/.well-known/host-meta.json', {
			method: 'GET',
		});
		const res = await app.request(req);
		const body = await res.json();
		expect(body).toMatchInlineSnapshot(`
			{
			  "links": [
			    {
			      "rel": "lrdd",
			      "template": "http://localhost/.well-known/webfinger?resource={uri}",
			    },
			  ],
			}
		`);
	});
});
