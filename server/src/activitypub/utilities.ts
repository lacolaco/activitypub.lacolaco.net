import { Activity } from './activity';
import { signRequest } from './signature';

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

export class ActivityPubAgent {
  constructor(readonly privateKey: string) {}

  async postActivity(url: URL, actorID: string, activity: Activity) {
    console.debug(`postActivity: ${url.toString()}`);
    console.debug(JSON.stringify(activity));
    const req = new Request(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/activity+json',
      },
      body: JSON.stringify(activity),
    });
    const signedReq = signRequest(req, actorID, this.privateKey);
    console.debug(JSON.stringify(Object.fromEntries(signedReq.headers.entries())));
    const res = await fetch(new Request(signedReq));
    if (!res.ok) {
      throw new Error(`postActivity: ${res.status} ${res.statusText}`);
    }
    console.debug(`postActivity: ${res.status} ${res.statusText}`);
    return res;
  }
}
