import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  app.get('/.well-known/host-meta', handleHostMetaXML);
  app.get('/.well-known/host-meta.json', handleHostMetaJSON);
};

const handleHostMetaXML: Handler = async (c) => {
  const accept = c.req.headers.get('Accept');
  const origin = c.get('origin');

  if (accept === 'application/json') {
    return c.redirect(`/.well-known/host-meta.json`);
  }

  const body = `<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
    <Link rel="lrdd" template="${origin}/.well-known/webfinger?resource={uri}"/>
</XRD>`;

  const res = c.text(body);
  res.headers.set('Content-Type', 'application/xrd+xml; charset=utf-8');
  return res;
};

const handleHostMetaJSON: Handler = async (c) => {
  const origin = c.get('origin');
  const body = {
    links: [
      {
        rel: 'lrdd',
        template: `${origin}/.well-known/webfinger?resource={uri}`,
      },
    ],
  };

  const res = c.json(body);
  return res;
};
