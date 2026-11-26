import { MiddlewareHandler } from 'hono';

export const assertAcceptHeader =
  (allowed: string[]): MiddlewareHandler =>
  async (c, next) => {
    const accept = c.req.header('Accept');
    if (accept == null) {
      return c.text('Bad Request', 400);
    }
    if (!allowed.some((a) => accept.includes(a))) {
      return c.text('Bad Request', 400);
    }
    await next();
  };

export const assertContentTypeHeader =
  (allowed: string[]): MiddlewareHandler =>
  async (c, next) => {
    const contentType = c.req.header('Content-Type');
    if (contentType == null) {
      return c.text('Bad Request', 400);
    }
    if (!allowed.some((a) => contentType.includes(a))) {
      return c.text('Bad Request', 400);
    }
    await next();
  };
