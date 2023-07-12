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
  if (typeof entity === 'object') {
    if ('id' in entity && typeof entity.id === 'string') {
      return new URL(entity.id);
    }
  }
  return null;
};

export class ActivityPubAgent {
  constructor(readonly privateKey: string) {}

  async postActivity(url: URL, actorID: string, activity: Activity) {
    const req = new Request(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/activity+json',
      },
      body: JSON.stringify(activity),
    });
    const signedReq = signRequest(req, actorID, this.privateKey);
    const res = await fetch(new Request(signedReq));
    if (!res.ok) {
      throw new Error(`postActivity: ${res.status} ${res.statusText}`);
    }
    return res;
  }
}
