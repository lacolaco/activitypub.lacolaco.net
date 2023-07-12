import { Hono } from 'hono';
import { logger } from 'hono/logger';
import { poweredBy } from 'hono/powered-by';

import { getPublicKey } from '@app/util/crypto';
import useActivityPub from '@app/web/ap';
import useHostMeta from '@app/web/host-meta';
import useNodeinfo from '@app/web/nodeinfo';
import useWebfinger from '@app/web/webfinger';
import { getConfigWithEnv } from './domain/config';
import { setupTracing } from './tracing';
import { AppContext } from './web/context';

async function createApplication(): Promise<Hono<AppContext>> {
  const app = new Hono<AppContext>();

  const config = await getConfigWithEnv();

  setupTracing();

  app.use('*', logger());
  app.use('*', poweredBy());
  app.use('*', async (c, next) => {
    c.set('Config', config);
    const privateKey = config.privateKeyPem;
    const publicKey = getPublicKey(privateKey);
    c.set('rsaKeyPair', { privateKey, publicKey });
    await next();
  });

  useNodeinfo(app);
  useHostMeta(app);
  useWebfinger(app);
  useActivityPub(app);

  return app;
}

export default createApplication;
