import { Hono } from 'hono';
import { logger } from 'hono/logger';
import { poweredBy } from 'hono/powered-by';

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

export default app;
