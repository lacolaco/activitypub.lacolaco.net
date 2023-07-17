import { Hono } from 'hono';
import { cors } from 'hono/cors';

import useAdmin from '@app/web/admin';
import useActivityPub from '@app/web/ap';
import useHostMeta from '@app/web/host-meta';
import useNodeinfo from '@app/web/nodeinfo';
import useWebfinger from '@app/web/webfinger';
import { Config } from './domain/config';
import { createLogger } from './logger';
import { withTracing } from './tracing';
import { AppContext, withConfig, withLogger, withOrigin } from './web/context';

async function createApplication(config: Config) {
  const app = new Hono<AppContext>();
  const logger = createLogger(config);

  app.use('*', withConfig(config));
  app.use('*', withOrigin());
  app.use('*', withTracing());
  app.use('*', withLogger(config, logger));

  console.log('config.clientOrigins', config.clientOrigins);
  app.use(
    '*',
    cors({
      origin: config.clientOrigins,
      credentials: true,
      allowHeaders: ['Content-Type', 'Authorization'],
    }),
  );

  useNodeinfo(app);
  useHostMeta(app);
  useWebfinger(app);
  useActivityPub(app);
  useAdmin(app, config);

  if (!config.isRunningOnCloud) {
    app.showRoutes();
  }

  return app;
}

export default createApplication;
