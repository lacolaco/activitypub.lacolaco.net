import * as crypto from 'node:crypto';

/**
 *
 * @see https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/importKey#pkcs_8_import
 */
export async function parsePrivateKey(pem: string) {
  const privateKey = crypto.createPrivateKey({
    key: pem,
    format: 'pem',
    type: 'pkcs8',
  });
  const publicKeyObject = crypto.createPublicKey(privateKey);
  const publicKeyPem = publicKeyObject
    .export({
      format: 'pem',
      type: 'spki',
    })
    .toString();

  return {
    privateKey,
    publicKeyPem,
  };
}
