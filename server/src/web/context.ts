import { getOrigin } from '@app/util/url';
import { MiddlewareHandler } from 'hono';

export type AppContext = {
  Variables: {
    readonly origin: string;
    readonly rsaKeyPair: {
      readonly publicKey: string;
      readonly privateKey: string;
    };
  };
};

export const withOrigin = (): MiddlewareHandler<AppContext> => async (c, next) => {
  const host = c.req.header('host');
  if (!host) {
    // host MUST be set in HTTP/1.1
    c.status(400);
    return c.json({ error: 'Bad Request' });
  }
  c.set('origin', getOrigin(host));
  await next();
};
