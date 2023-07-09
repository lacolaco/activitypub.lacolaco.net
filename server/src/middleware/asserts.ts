import { MiddlewareHandler } from 'hono';

export const assertAcceptHeader =
	(allowed: string[]): MiddlewareHandler =>
	async (c, next) => {
		const accept = c.req.headers.get('Accept');
		if (accept == null) {
			c.status(400);
			return c.text('Bad Request');
		}
		if (!allowed.some((a) => accept.includes(a))) {
			c.status(400);
			return c.text('Bad Request');
		}
		await next();
	};

export const assertContentTypeHeader =
	(allowed: string[]): MiddlewareHandler =>
	async (c, next) => {
		const contentType = c.req.headers.get('Content-Type');
		if (contentType == null) {
			c.status(400);
			return c.text('Bad Request');
		}
		if (!allowed.some((a) => contentType.includes(a))) {
			c.status(400);
			return c.text('Bad Request');
		}
		await next();
	};
