import { serve } from '@hono/node-server';
import createApplication from './app';

const port = Number(process.env.PORT || 8080);

createApplication()
  .then((app) => {
    serve({ fetch: app.fetch, port }).once('listening', () => {
      console.log(`Listening on http://localhost:${port}`);
    });
  })
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
