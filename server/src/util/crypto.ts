import { createPublicKey } from 'node:crypto';

export function getPublicKey(privateKey: string) {
  const publicKey = createPublicKey(privateKey);
  return publicKey.export({ type: 'spki', format: 'pem' }).toString();
}
