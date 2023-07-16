import { Config } from '@app/domain/config';
import { Logger, createLoggerWithTrace } from '@app/logger';
import { getOrigin } from '@app/util/url';
import { context, trace } from '@opentelemetry/api';
import { MiddlewareHandler } from 'hono';

export type AppContext = {
  Variables: {
    readonly origin: string;
    readonly config: Config;
    readonly logger: Logger;
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

export function withLogger(base: Logger): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    const spanContext = trace.getSpanContext(context.active());
    const logger = spanContext ? createLoggerWithTrace(base, spanContext) : base;
    c.set('logger', logger);

    const { method, url } = c.req;
    logger.info(`--> ${method} ${url}`);
    await next();
    logger.info(`<-- ${method} ${url} ${c.res.status}`);
  };
}
