import { AP } from '@activity-kit/types';
import { Activity } from './activity';
import { signRequest } from './signature';

export const getEntityID = (entity?: undefined | null | AP.EntityReference | AP.EntityReference[]) => {
  if (!entity || Array.isArray(entity)) {
    return null;
  }

  if (entity instanceof URL) {
    return entity;
  }

  if ('id' in entity) {
    return entity.id ?? null;
  }

  return null;
};

export class ActivityPubAgent {
  constructor(readonly privateKey: string) {}

  postActivity(url: URL, actorID: string, activity: Activity) {
    const req = new Request(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/activity+json',
      },
      body: JSON.stringify(activity),
    });
    const signedReq = signRequest(req, actorID, this.privateKey);
    return fetch(new Request(signedReq));
  }
}
