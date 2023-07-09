import { Hono } from 'hono';
import { logger } from 'hono/logger';
import { poweredBy } from 'hono/powered-by';
import { serve } from '@hono/node-server';
import { serveStatic } from '@hono/node-server/serve-static';

import useActivityPub from '@app/web/ap';
import useHostMeta from '@app/web/host-meta';
import useNodeinfo from '@app/web/nodeinfo';
import useWebfinger from '@app/web/webfinger';

const app = new Hono();

app.use('*', logger());
app.use('*', poweredBy());

useNodeinfo(app);
useHostMeta(app);
useWebfinger(app);
useActivityPub(app);

app.use('/*', serveStatic({ root: './' }));

export default app;
