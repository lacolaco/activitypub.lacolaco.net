import { getConfigWithEnv } from '@app/domain/config';
import { describe, expect, test } from 'vitest';
import { createDigest, parseSignatureString, signHeaders, verifySignature } from './signature';

describe('signature', () => {
  describe('createDigest', () => {
    test('creates a digest', async () => {
      const body = { hello: 'world' };
      console.log(JSON.stringify(body));
      const digest = createDigest(body);
      // echo -n '{"hello":"world"}' | sha256sum | xxd -r -p | base64
      expect(digest).toBe('k6I5cakU5erL8KjSUVTNownDwccvu5kU1Hxg88toFYg=');
    });
  });

  describe('parseSignatureString', () => {
    test('parses a signature string', async () => {
      const signature =
        'keyId="https://example.com/users/1",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="abc"';
      const parsed = parseSignatureString(signature);
      expect(parsed).toEqual({
        keyId: 'https://example.com/users/1',
        algorithm: 'rsa-sha256',
        headers: '(request-target) host date digest',
        signature: 'abc',
      });
    });
  });

  describe('signHeaders', () => {
    test('signs a request', async () => {
      const privateKey = (await getConfigWithEnv()).privateKey;

      const inbox = 'https://remote.example.com/inbox';
      const body = { hello: 'world' };
      const actorID = 'https://example.com/users/1';

      const headers = await signHeaders('POST', inbox, body, actorID, privateKey, new Date('2021-01-01T00:00:00Z'));

      const signature = headers['Signature'];
      expect(signature).not.toBe(null);
    });
  });

  describe('verifySignature', () => {
    test('verifies a request from minidon to misskey', async () => {
      const inbox = new URL('https://misskey.io/users/9bdgn9zxoi/inbox');
      const body = {
        '@context': 'https://www.w3.org/ns/activitystreams',
        id: 'https://minidon-debug-alice.lacolaco.workers.dev/u/alice/s/21029dbb-4d12-4c0a-9ca7-54a26fbbb7fc',
        type: 'Accept',
        actor: 'https://minidon-debug-alice.lacolaco.workers.dev/u/alice',
        object: {
          '@context': [
            'https://www.w3.org/ns/activitystreams',
            'https://w3id.org/security/v1',
            {
              manuallyApprovesFollowers: 'as:manuallyApprovesFollowers',
              sensitive: 'as:sensitive',
              Hashtag: 'as:Hashtag',
              quoteUrl: 'as:quoteUrl',
              toot: 'http://joinmastodon.org/ns#',
              Emoji: 'toot:Emoji',
              featured: 'toot:featured',
              discoverable: 'toot:discoverable',
              schema: 'http://schema.org#',
              PropertyValue: 'schema:PropertyValue',
              value: 'schema:value',
              misskey: 'https://misskey-hub.net/ns#',
              _misskey_content: 'misskey:_misskey_content',
              _misskey_quote: 'misskey:_misskey_quote',
              _misskey_reaction: 'misskey:_misskey_reaction',
              _misskey_votes: 'misskey:_misskey_votes',
              isCat: 'misskey:isCat',
              vcard: 'http://www.w3.org/2006/vcard/ns#',
            },
          ],
          id: 'https://misskey.io/follows/9h4ovb4nmt',
          type: 'Follow',
          actor: 'https://misskey.io/users/9bdgn9zxoi',
          object: 'https://minidon-debug-alice.lacolaco.workers.dev/u/alice',
        },
      };

      const req = new Request(inbox, {
        method: 'POST',
        headers: {
          Host: 'misskey.io',
          Date: 'Thu, 13 Jul 2023 12:26:59 GMT',
          Digest: 'SHA-256=pK20diGAwwwlT3/kZsXnHGYDX1FEPDqTi4htwA81fcA=',
          Signature:
            'keyId="https://minidon-debug-alice.lacolaco.workers.dev/u/alice",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="VCyeH2WRX2D3i8c8zCxFYx7q9T8U8ZuD150fAn3O39oTZSj1n5LufEr/PU+DdWO3uAy+zZcGeeJd0ksCVzbQyOjJhuwN32zy2HOjFULLmm2rjA3cZbAnAn6bMEloy9MzxDEvLMdkF/vCIoPLIOtooSMM86S4O8suXvuwXPi9MbV2b+DPunNv4RjZxTSprlV4w2b/XM+IVFFbLcqpRk33QMH3rKXe+XEXPc/SoaCb3A0p+G7hY72Sfqtwt03aE6MOK+56PWZDftQn8trYxC+hcaj+ii0UZrU3QQSow9Y/La2PEg5Su9XTSsxFsbaM45oV+tgrrchuMP4XJDRB3uZa70zKeX6OcqU0csK1zfI9w7sr4fY+eun4YE2kDt8spG8tp8qk8pXtU1z39Kjy24bOleq5I1LgU0Ua9OMnyoTQjnsuupdK6lxeblzZmxPvC3oEzouEDH7EH3NeSNie5jkRQeya9U3cbBgYThYP+u+RMqO7sHhjUSmwU7653XV8cHBMd7PxHJFq4sfxvmnvtryHw9cpzkAST2ZgupvL9RKB4WdK1XaMrXv+OVxiyk4L+A/m/pe35oN1R01Q6D+21dlnIB6HlikkHU/R3ycdll6mmqrGRt/jrKn+4GWuukgclq9N2q4Zvmz7V+nxz6ezk/M8sV8KwM4eV9lnv3Af8FP5cv0="',
          Accept: 'application/activity+json',
          'Content-Type': 'application/activity+json',
          'Accept-Encoding': 'gzip',
          'User-Agent': 'Minidon/0.0.0 (+https://minidon-debug-alice.lacolaco.workers.dev/)',
        },
        body: JSON.stringify(body),
      });

      const ok = await verifySignature(req, async () => ({
        id: 'https://minidon-debug-alice.lacolaco.workers.dev/u/alice',
        type: 'Key',
        owner: 'https://minidon-debug-alice.lacolaco.workers.dev/u/alice',
        publicKeyPem:
          '-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAu82gXsBVXXJTfqdB9ccL\nRCGtvYF7Xj4Jb8GkdokABiaQpMfv49nCfSpCmTKIEyaxlnStDasdNaK8pkUTUmjF\noUzMD4IrsDYFqHoO7I1q/9oBUn9OkxcpNl3fTXfTOjPuoqtdqe1Qck/9X5ptRcX4\nmMiO/DH/pJNUV0zwHvjH51D0YWD7N1/Mkkc/2O3FwCLHeFanxxdOa1MMLKU+zG4v\n/OEKmSNg6M0avLUPk9EKe0rYPltrM8q+dAg37r4FMf88CCDJexpSv+ix/KsXy7HO\ngAzvnY8ptb+COeGSIoL1BGFMeRf+c9DGCjwPJ0sqkcfB0ravZGt4fRg2RsHgj57/\necpnRxnyS8Zcir12Y8QQ3DEr7jA+QfHKUooCqSOz26Q7bNJdk/0Ay35DpPrbCrWN\nqobQV3mUtcMWEE6xaWHm+fmLMxXSdbfcXvgH93yjUqetqLQYADzQlIClYvo94lYs\n0jPkjRRyXqt9bqk2rJiGMFDPG2dgh8V0DAbL4DfCDtFGkzdzbMFvHmWgyIa0TaO2\nwRZXCtfn9+TItD7DZziEjmsfC5NnTe2dQPM9Sk8qb44P179GqbGO8PHbETpq79X7\nkOvE6TUkZMiSvlw59UuH9uQxST93R1YWNrmeBUEw4I7SpZwpRDwxl9MeMt5MHGld\nMdBi6Jr00O+JlnRpOMaAe4ECAwEAAQ==\n-----END PUBLIC KEY-----\n',
      }));

      expect(ok).toBe(true);
    });

    test('verifies a request signed by me', async () => {
      const config = await getConfigWithEnv();
      const privateKey = config.privateKey;
      const publicKeyPem = config.publicKeyPem;

      const inbox = 'https://remote.example.com/inbox';
      const body = { hello: 'world' };
      const actorID = 'https://example.com/users/1';

      const headers = await signHeaders('POST', inbox, body, actorID, privateKey);
      console.log(headers);

      const req = new Request(inbox, {
        method: 'POST',
        headers,
        body: JSON.stringify(body),
      });

      const ok = await verifySignature(req, async () => ({
        id: 'https://example.com/users/1#key',
        owner: 'https://example.com/users/1',
        publicKeyPem,
      }));

      expect(ok).toBe(true);
    });
  });
});
