import { Config } from '@app/domain/config';
import { Logger } from '@app/logger';
import { getOrigin } from '@app/util/url';
import { MiddlewareHandler } from 'hono';

export type AppContext = {
  Variables: {
    readonly origin: string;
    readonly config: Config;
  };
};

export function withOrigin(): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    const host = c.req.header('host');
    if (!host) {
      c.status(400);
      return c.json({ error: 'Bad Request' });
    }
    c.set('origin', getOrigin(host));
    await next();
  };
}

export function withConfig(config: Config): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    c.set('config', config);
    await next();
  };
}

export function withLogger(logger: Logger): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    const { method, url } = c.req;
    logger.info(`--> ${method} ${url}`);
    await next();
    logger.info(`<-- ${method} ${url} ${c.res.status}`);
  };
}
