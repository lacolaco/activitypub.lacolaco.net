import parser, { Sha256Signer } from 'activitypub-http-signatures';

export function getPublicKeyID(actorID: string): string {
  return `${actorID}#main-key`;
}

/**
 * Public key for HTTP Signatures.
 *
 * @see https://docs.joinmastodon.org/spec/activitypub/#publicKey
 */
export type PublicKey = {
  id: string;
  owner: string;
  publicKeyPem: string;
};

export function signRequest(req: Request, actorID: string, privateKey: string) {
  const { url, method, headers } = req;
  const headerNames = Object.keys(headers);
  const headersObject = Object.fromEntries(headers.entries());
  const publicKeyId = getPublicKeyID(actorID);

  const signer = new Sha256Signer({ publicKeyId, privateKey, headerNames });
  const signature = signer.sign({ url, method, headers: headersObject });

  req.headers.set('Signature', signature);

  return req;
}

export async function verifySignature(req: { url: string; method: string; headers: Headers }) {
  const { url, method, headers } = req;
  const headersObject = Object.fromEntries(headers.entries());
  const signature = parser.parse({ url, method, headers: headersObject });

  const publicKey = await fetchPublicKey(signature.keyId);
  const success = signature.verify(publicKey.publicKeyPem);

  return success;
}

async function fetchPublicKey(keyID: string) {
  const res = await fetch(keyID, {
    headers: {
      accept: 'application/ld+json, application/json',
    },
  });
  const { publicKey } = (await res.json()) as { publicKey: PublicKey };
  return publicKey;
}
