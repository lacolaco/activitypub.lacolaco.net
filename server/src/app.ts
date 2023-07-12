import { Hono } from 'hono';
import { logger } from 'hono/logger';

import { getPublicKey } from '@app/util/crypto';
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
  app.use('*', withOrigin());
  app.use('*', async (c, next) => {
    const privateKey = config.privateKeyPem;
    const publicKey = getPublicKey(privateKey);
    c.set('rsaKeyPair', { privateKey, publicKey });
    await next();
  });

  useNodeinfo(app);
  useHostMeta(app);
  useWebfinger(app);
  useActivityPub(app);

  if (isDevelopment) {
    app.routes.forEach((route) => {
      console.log(`${route.method} ${route.path}`);
    });
  }

  return app;
}

export default createApplication;
