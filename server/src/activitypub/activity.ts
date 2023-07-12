import { AP, isTypeOf } from '@activity-kit/types';
import { contextURIs } from './context';

export type Activity = AP.Activity;
export type FollowActivity = AP.Follow;

export function getActorOf(activity: Activity): AP.Actor {
  const actor = activity.actor;
  if (isTypeOf<AP.Actor>(actor, AP.ActorTypes)) {
    return actor;
  }
  throw new Error('invalid actor');
}

export function isFollowActivity(activity: Activity): activity is AP.Follow {
  return activity.type === 'Follow';
}

export function buildAcceptAcivity(actorID: URL, object: Activity) {
  const id = new URL(`${actorID}/accept/${object.id}`);
  return {
    '@context': contextURIs,
    type: 'Accept',
    id,
    actor: actorID,
    object: object,
  } satisfies AP.Accept;
}
