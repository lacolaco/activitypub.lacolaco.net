import { getTracer } from '@app/tracing';
import { KeyObject, createHash, sign, verify } from 'node:crypto';
import { Person } from './person';
import { getID } from './utilities';

export function getPublicKeyID(actorID: string | URL): string {
  return `${actorID.toString()}#key`;
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

export function withPublicKey(person: Person, publicKey: string): Person {
  const id = getID(person);
  if (id == null) {
    throw new Error('person.id is null');
  }

  return {
    ...person,
    publicKey: {
      id: getPublicKeyID(id.toString()),
      owner: id.toString(),
      publicKeyPem: publicKey,
    },
  };
}

/**
 * Convert string to bytes.
 */
function stob(s: string) {
  return Uint8Array.from(s, (c) => c.charCodeAt(0));
}

/**
 * Convert bytes to string.
 */
function btos(b: ArrayBuffer) {
  return String.fromCharCode(...new Uint8Array(b));
}

function createSignString(req: { method: string; url: URL }, headers: Record<string, string>, headerNames?: string[]) {
  headerNames = headerNames || Object.keys(headers);
  return [
    `(request-target): ${req.method.toLowerCase()} ${req.url.pathname}`,
    ...headerNames.map((name) => `${name.toLowerCase()}: ${headers[name]}`),
  ].join('\n');
}

function trimQuotes(s: string) {
  if (s.startsWith('"')) s = s.slice(1);
  if (s.endsWith('"')) s = s.slice(0, -1);
  return s;
}

export function parseSignatureString(signature: string): Record<string, string> {
  const parts = signature.split(',');
  const map = Object.fromEntries(parts.map((part) => part.split('=').map((s) => trimQuotes(s))));
  return map;
}

export function createDigest(body: object) {
  const hash = createHash('SHA-256');
  hash.write(JSON.stringify(body));
  hash.end();
  return hash.digest('base64');
}

export async function signHeaders(
  method: 'POST',
  url: URL,
  body: object,
  publicKeyID: string,
  privateKey: KeyObject,
  now = new Date(),
) {
  const dateStr = now.toUTCString();
  const digest = createDigest(body);

  const signString = createSignString(
    { method, url },
    {
      host: url.host,
      date: dateStr,
      digest: `SHA-256=${digest}`,
    },
  );
  const signature = sign('SHA-256', Buffer.from(signString), privateKey).toString('base64');

  const headers = {
    Host: url.host,
    Date: dateStr,
    Digest: `SHA-256=${digest}`,
    Signature:
      `keyId="${publicKeyID}",` +
      `algorithm="rsa-sha256",` +
      `headers="(request-target) host date digest",` +
      `signature="${signature}"`,
  };
  return headers;
}

export type ResolvePublicKeyFn = (keyID: string) => Promise<PublicKey>;

export async function verifySignature(
  req: { url: string; method: string; headers: Headers },
  resolvePublicKey: ResolvePublicKeyFn = fetchPublicKey,
) {
  const { url, method, headers } = req;
  const sigHeader = headers.get('Signature');
  if (sigHeader == null) {
    throw new Error('Signature header is missing');
  }
  const signatureFields = parseSignatureString(sigHeader);
  const signature = signatureFields.signature;
  const headerNames = signatureFields.headers.split(/\s+/) ?? ['(request-target)', 'host', 'date', 'digest'];
  const keyID = signatureFields.keyId;
  if (keyID == null) {
    throw new Error('keyId is missing');
  }
  const { publicKeyPem } = await resolvePublicKey(keyID);

  const headersObject = Object.fromEntries(headers.entries());
  const signString = createSignString({ method, url: new URL(url) }, headersObject, headerNames);

  try {
    verify('SHA-256', Buffer.from(signString), publicKeyPem, Buffer.from(signature));
    return true;
  } catch (err) {
    console.error(err);
    throw new Error('Signature verification failed');
  }
}

async function fetchPublicKey(keyID: string) {
  return getTracer().startActiveSpan('fetchPublicKey', async (span) => {
    span.setAttribute('keyID', keyID);
    console.debug(`fetchPublicKey: ${keyID}`);

    const res = await fetch(keyID, {
      headers: {
        accept: 'application/activity+json',
      },
    });
    if (!res.ok) {
      throw new Error(`fetchPublicKey: ${res.status} ${res.statusText}`);
    }
    const body = (await res.json()) as { publicKey: PublicKey };
    console.debug(JSON.stringify(body));
    return body.publicKey;
  });
}
