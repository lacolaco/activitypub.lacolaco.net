import { serve } from '@hono/node-server';
import createApplication from './app';
import { setupTracing } from './tracing';
import { getConfigWithEnv } from './domain/config';

const port = Number(process.env.PORT || 8080);

const config = getConfigWithEnv();
setupTracing(config);

createApplication(config)
  .then((app) => {
    serve({ fetch: app.fetch, port }).once('listening', () => {
      console.log(`Listening on http://localhost:${port}`);
    });
  })
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
