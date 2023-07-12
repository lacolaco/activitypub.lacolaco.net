import { generateKeyPairSync } from 'node:crypto';

export function setup() {
  console.log('test setup');

  const keys = generateKeyPairSync('rsa', {
    modulusLength: 2048,
    publicKeyEncoding: {
      type: 'spki',
      format: 'pem',
    },
    privateKeyEncoding: {
      type: 'pkcs8',
      format: 'pem',
    },
  });
  process.env['RSA_PRIVATE_KEY'] = keys.privateKey;
}

export function teardown() {}
