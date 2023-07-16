import { parsePrivateKey } from '@app/util/crypto';
import { KeyObject } from 'node:crypto';

export type Config = {
  readonly privateKey: KeyObject;
  readonly publicKeyPem: string;
  readonly clientOrigins: string[];
  readonly isRunningOnCloud: boolean;
  readonly gcpProjectID: string;
};

export function getConfigWithEnv(): Config {
  const privateKeyPem = process.env['RSA_PRIVATE_KEY'];
  if (privateKeyPem == null) {
    throw new Error('RSA_PRIVATE_KEY is not set');
  }
  const { privateKey, publicKeyPem } = parsePrivateKey(privateKeyPem);
  const clientOrigins = (process.env['CLIENT_ORIGIN'] ?? '').split(',');
  const isRunningOnCloud = isRunningOnCloudRun();
  const gcpProjectID = getGCPProjectID();
  return {
    privateKey,
    publicKeyPem,
    clientOrigins,
    isRunningOnCloud,
    gcpProjectID,
  };
}

function isRunningOnCloudRun(): boolean {
  return process.env['K_SERVICE'] !== undefined;
}

export function getGCPProjectID() {
  return process.env['GCP_PROJECT'] ?? '';
}
