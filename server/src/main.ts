import { serve } from '@hono/node-server';
import app from './app';

app.routes.forEach((route) => {
  console.log(`${route.method} ${route.path}`);
});

const port = Number(process.env.PORT || 8080);

serve({ fetch: app.fetch, port }).once('listening', () => {
  console.log(`Listening on http://localhost:${port}`);
});
