import { KeyObject } from 'node:crypto';
import { signHeaders } from './signature';
import { AnyActivity, ObjectOrURI, URI } from './schema';

export const getURI = (object: ObjectOrURI) => {
  if (typeof object === 'string') {
    return object;
  }
  const { id } = object;
  if (typeof id === 'string') {
    return id;
  }
  return null;
};

export async function postActivity(inbox: URI, activity: AnyActivity, publicKeyID: string, privateKey: KeyObject) {
  console.debug(`postActivity: ${inbox.toString()}`);
  console.debug(JSON.stringify(activity));
  const headers = await signHeaders('POST', inbox, activity, publicKeyID, privateKey);
  const res = await fetch(inbox, {
    method: 'POST',
    headers: {
      ...headers,
      Accept: 'application/activity+json',
      'Content-Type': 'application/activity+json',
      'Accept-Encoding': 'gzip',
      'User-Agent': `activitypub.lacolaco.net/1.0`,
    },
    body: JSON.stringify(activity),
  });
  if (!res.ok) {
    throw new Error(`postActivity: ${res.status} ${res.statusText}`);
  }
  console.debug(`postActivity: ${res.status} ${res.statusText}`);
  return res;
}
