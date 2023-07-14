import { KeyObject } from 'crypto';
import { Activity } from './activity';
import { signHeaders } from './signature';

export const getID = (entity: unknown) => {
  if (entity == null) {
    return null;
  }
  if (typeof entity === 'string') {
    return new URL(entity);
  }
  if (entity instanceof URL) {
    return entity;
  }
  if (typeof entity === 'object' && 'id' in entity) {
    if (typeof entity.id === 'string') {
      return new URL(entity.id);
    }
    if (entity.id instanceof URL) {
      return entity.id;
    }
  }
  return null;
};

export async function postActivity(inbox: URL, activity: Activity, publicKeyID: string, privateKey: KeyObject) {
  console.debug(`postActivity: ${inbox.toString()}`);
  console.debug(JSON.stringify(activity));
  const headers = await signHeaders('POST', inbox, activity, publicKeyID, privateKey);
  console.debug(JSON.stringify(headers));
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
