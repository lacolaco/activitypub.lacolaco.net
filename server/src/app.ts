import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';

import useAdmin from '@app/web/admin';
import useActivityPub from '@app/web/ap';
import useHostMeta from '@app/web/host-meta';
import useNodeinfo from '@app/web/nodeinfo';
import useWebfinger from '@app/web/webfinger';
import { getConfigWithEnv } from './domain/config';
import { setupTracing } from './tracing';
import { AppContext, withOrigin } from './web/context';

async function createApplication(): Promise<Hono<AppContext>> {
  const app = new Hono<AppContext>();

  const config = await getConfigWithEnv();
  const isDevelopment = !config.isRunningOnCloud;

  setupTracing();

  app.use('*', logger());
  app.use(
    '*',
    cors({
      origin: config.clientOrigins,
      credentials: true,
      allowMethods: ['GET', 'POST', 'OPTIONS'],
    }),
  );
  app.use('*', withOrigin());
  app.use('*', async (c, next) => {
    c.set('config', config);
    await next();
  });

  useNodeinfo(app);
  useHostMeta(app);
  useWebfinger(app);
  useActivityPub(app);
  useAdmin(app);

  if (isDevelopment) {
    app.routes.forEach((route) => {
      console.log(`${route.method} ${route.path}`);
    });
  }

  return app;
}

export default createApplication;
