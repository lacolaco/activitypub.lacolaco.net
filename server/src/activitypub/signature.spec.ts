import { getConfigWithEnv } from '@app/domain/config';
import { getPublicKey } from '@app/util/crypto';
import { describe, expect, test } from 'vitest';
import { signRequest, verifySignature } from './signature';

describe('signature', () => {
  describe('signRequest', () => {
    test('signs a request', async () => {
      const privateKey = (await getConfigWithEnv()).privateKeyPem;

      const req = new Request('https://remote.example.com/inbox', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/activity+json',
        },
        body: JSON.stringify({ hello: 'world' }),
      });
      const actorID = 'https://example.com/users/1';

      const signed = signRequest(req, actorID, privateKey);

      const signature = signed.headers.get('Signature');
      expect(signature).not.toBe(null);
    });
  });

  describe('verifySignature', () => {
    test('verifies a request', async () => {
      const privateKey = (await getConfigWithEnv()).privateKeyPem;
      const publicKey = getPublicKey(privateKey);

      const req = new Request('https://remote.example.com/inbox', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/activity+json',
        },
        body: JSON.stringify({ hello: 'world' }),
      });

      const signed = signRequest(req, 'https://example.com/users/1', privateKey);

      const ok = await verifySignature(signed, async () => ({
        id: 'https://example.com/users/1#key',
        owner: 'https://example.com/users/1',
        publicKeyPem: publicKey,
      }));

      expect(ok).toBe(true);
    });
  });
});
