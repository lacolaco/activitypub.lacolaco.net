import { AP, isType } from '@activity-kit/types';
import { contextURIs } from './context';
import { getID } from './utilities';

export type Activity = AP.Activity & Record<string, unknown>;

export type FollowActivity = Activity & AP.Follow;

export function isFollowActivity(object: unknown): object is FollowActivity {
  return isType<AP.Follow>(object, AP.ActivityTypes.FOLLOW);
}

export type UndoActivity = Activity & AP.Undo;

export function isUndoActivity(object: unknown): object is UndoActivity {
  return isType<AP.Undo>(object, AP.ActivityTypes.UNDO);
}

export type AcceptActivity = Activity & AP.Accept;

export function isAcceptActivity(object: unknown): object is AcceptActivity {
  return isType<AP.Accept>(object, AP.ActivityTypes.ACCEPT);
}

export function buildAcceptAcivity(actorID: URL, object: Activity) {
  const id = new URL(`${actorID}/accept/${object.id}`);
  return {
    '@context': contextURIs,
    type: 'Accept',
    id,
    actor: actorID,
    object: object,
    to: [getID(object.actor)],
    published: new Date().toISOString(),
  } satisfies AcceptActivity;
}
