import parser, { Sha256Signer } from 'activitypub-http-signatures';
import { Person } from './person';
import { getID } from './utilities';
import { getTracer } from '@app/tracing';

export function getPublicKeyID(actorID: string): string {
  return `${actorID}#key`;
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

export function withPublicKey(entity: Person, publicKey: string): Person {
  const id = getID(entity);
  if (id == null) {
    throw new Error('person.id is null');
  }

  return {
    ...entity,
    publicKey: {
      id: getPublicKeyID(id.toString()),
      owner: id.toString(),
      publicKeyPem: publicKey,
    },
  };
}

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
  return getTracer().startActiveSpan('fetchPublicKey', async (span) => {
    span.setAttribute('keyID', keyID);

    const res = await fetch(keyID, {
      headers: {
        accept: 'application/activity+json',
      },
    });
    const { publicKey } = (await res.json()) as { publicKey: PublicKey };
    return publicKey;
  });
}
