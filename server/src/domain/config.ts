import { parsePrivateKey } from '@app/util/crypto';
import { KeyObject } from 'crypto';
import { GoogleAuth } from 'google-auth-library';

export type Config = {
  readonly privateKey: KeyObject;
  readonly publicKeyPem: string;
  readonly gcpProjectID: string;
  readonly clientOrigin: string;
  readonly isRunningOnCloud: boolean;
};

export async function getConfigWithEnv(): Promise<Config> {
  const privateKeyPem = process.env['RSA_PRIVATE_KEY'];
  if (privateKeyPem == null) {
    throw new Error('RSA_PRIVATE_KEY is not set');
  }
  const { privateKey, publicKeyPem } = await parsePrivateKey(privateKeyPem);
  const clientOrigin = process.env['CLIENT_ORIGIN'] ?? '';
  const googleCredentials = await findGoogleCredentials();
  const isRunningOnCloud = isRunningOnCloudRun();

  return {
    privateKey,
    publicKeyPem,
    gcpProjectID: googleCredentials.projectId ?? '',
    clientOrigin,
    isRunningOnCloud,
  };
}

function isRunningOnCloudRun(): boolean {
  return process.env['K_SERVICE'] !== undefined;
}

async function findGoogleCredentials() {
  try {
    const credentials = await new GoogleAuth().getApplicationDefault();
    return credentials;
  } catch (e) {
    throw new Error('Google credentials not found');
  }
}
