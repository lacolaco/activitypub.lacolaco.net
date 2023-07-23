import { AppContext } from '@app/web/context';
import { MiddlewareHandler } from 'hono';
import { verify } from 'jsonwebtoken';

const jwtKeysURL = 'https://lacolaco.cloudflareaccess.com/cdn-cgi/access/certs';
const aud = '20c040a974a2b0878ea9cf37df82bb1a5d7680ce262c27088eae68901ba4888d';

export function verifyJWT(): MiddlewareHandler<AppContext> {
  return async (c, next) => {
    const token = c.req.header('Cf-Access-Jwt-Assertion');
    if (!token) {
      c.status(401);
      return c.json({ error: 'Unauthorized' });
    }
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
