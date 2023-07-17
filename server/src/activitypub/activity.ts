import { randomUUID } from 'node:crypto';
import { contextURIs, contextURIsWithExtensions } from './context';
import {
  AcceptActivity,
  ActivityPubObject,
  ActivityStreamsObject,
  AnyActivity,
  FollowActivity,
  ObjectRef,
  URI,
  UndoActivity,
} from './schema';

export function isFollowActivity(object: ObjectRef): object is FollowActivity {
  if (typeof object === 'string') {
    return false;
  }
  return object.type === 'Follow';
}

export function isUndoActivity(object: ObjectRef): object is UndoActivity {
  if (typeof object === 'string') {
    return false;
  }
  return object.type === 'Undo';
}

export function isAcceptActivity(object: ObjectRef): object is AcceptActivity {
  if (typeof object === 'string') {
    return false;
  }
  return object.type === 'Accept';
}

export function buildAcceptAcivity(actorID: URI, object: AnyActivity, acceptID: string = randomUUID()) {
  const id = new URL(`${actorID}/accept/${acceptID}`);
  return AcceptActivity.parse({
    '@context': contextURIs,
    type: 'Accept',
    id,
    actor: actorID,
    object: object,
  });
}

export function buildCreateActivity(actorID: URI, object: ActivityPubObject, createID: string = randomUUID()) {
  const id = new URL(`${actorID}/create/${createID}`);
  return AcceptActivity.parse({
    '@context': contextURIsWithExtensions,
    type: 'Create',
    id,
    actor: actorID,
    object: object,
  });
}
