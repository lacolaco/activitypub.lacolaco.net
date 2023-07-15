import { searchPerson } from '@app/usecase/admin/search-person';
import { Handler, Hono, MiddlewareHandler } from 'hono';
import { verify } from 'jsonwebtoken';
import { AppContext } from '../context';

const jwtKeysURL = 'https://lacolaco.cloudflareaccess.com/cdn-cgi/access/certs';
const aud = '20c040a974a2b0878ea9cf37df82bb1a5d7680ce262c27088eae68901ba4888d';

function verifyJWT(): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    const config = c.get('config');
    if (!config.isRunningOnCloud) {
      await next();
      return;
    }

    const token = c.req.header('Cf-Access-Jwt-Assertion');
    if (!token) {
      c.status(401);
      return c.json({ error: 'Unauthorized' });
    }
    console.log(token);
    const { public_certs } = (await fetch(jwtKeysURL).then((res) => res.json())) as {
      public_certs: { kid: string; cert: string }[];
    };

    try {
      await new Promise((resolve, reject) => {
        verify(
          token,
          ({ kid }, callback) => {
            const key = public_certs.find((k) => k.kid === kid);
            if (key == null) {
              callback(new Error('Invalid kid'));
              return;
            }
            callback(null, key.cert);
          },
          { audience: aud },
          (err, decoded) => {
            err ? reject(err) : resolve(decoded);
          },
        );
      });
    } catch (e) {
      console.error(e);
      c.status(401);
      return c.json({ error: 'Unauthorized' });
    }

    await next();
  };
}

export default (app: Hono<AppContext>) => {
  const adminRoutes = new Hono<AppContext>();

  adminRoutes.use('*', verifyJWT());
  adminRoutes.get('/users/show/:username', (c) => {
    const headers = Object.fromEntries(c.req.headers.entries());
    console.log(JSON.stringify(headers));

    c.status(404);
    return c.json({ error: 'Not Found' });
  });
  // apiRoutes.get('/search/person/:resource', handleSearchPerson);

  app.route('/admin', adminRoutes);
};

const handleSearchPerson: Handler<AppContext> = async (c) => {
  const resource = c.req.param('resource');

  try {
    const person = await searchPerson(resource);
    return c.json(person);
  } catch (e) {
    c.status(404);
    return c.json({ error: 'Not Found' });
  }
};
