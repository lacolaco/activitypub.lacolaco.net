import { describe, expect, test } from 'vitest';
import { parsePrivateKey } from './crypto';

describe('crypto', () => {
  describe('parsePrivateKey', () => {
    test('should parse private key string', async () => {
      const pem = process.env.RSA_PRIVATE_KEY!;
      const { privateKey, publicKeyPem } = await parsePrivateKey(pem);
      expect(privateKey).toBeDefined();
      expect(publicKeyPem).toBeDefined();
    });
  });
});
