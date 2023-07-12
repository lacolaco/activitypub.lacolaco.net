import { AP, isType } from '@activity-kit/types';
import { contextURIs } from './context';

export type Activity = AP.Activity;

export function getActorOf(activity: Activity): AP.Actor {
  const actor = activity.actor;
  if (isType<AP.Actor>(actor, AP.ActorTypes.PERSON)) {
    return actor;
  }
  throw new Error(`invalid actor: ${JSON.stringify(actor)}`);
}

export type FollowActivity = AP.Follow;

export function isFollowActivity(object: unknown): object is FollowActivity {
  return isType<AP.Follow>(object, AP.ActivityTypes.FOLLOW);
}

export type UndoActivity = AP.Undo;

export function isUndoActivity(object: unknown): object is UndoActivity {
  return isType<AP.Undo>(object, AP.ActivityTypes.UNDO);
}

export type AcceptActivity = AP.Accept;

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
  } satisfies AcceptActivity;
}
