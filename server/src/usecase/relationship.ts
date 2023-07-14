import { Config } from '@app/domain/config';
import { RemoteUser } from '@app/domain/remote-user';
import { User } from '@app/domain/user';
import { UserFollowersRepository } from '@app/repository/user-followers';
import { AppContext } from '@app/web/context';
import { Context } from 'hono';
import { FollowActivity, UndoActivity, buildAcceptAcivity } from '../activitypub/activity';
import { buildPerson, fetchPersonByID } from '../activitypub/person';
import { getPublicKeyID } from '../activitypub/signature';
import { getID, postActivity } from '../activitypub/utilities';

export async function getUserFollowers(user: User): Promise<RemoteUser[]> {
  const followersRepo = new UserFollowersRepository();
  const followers = await followersRepo.list(user);
  return followers;
}

export async function acceptFollowRequest(
  config: Config,
  c: Context<AppContext>,
  user: User,
  activity: FollowActivity,
) {
  const actorID = getID(activity.actor);
  if (actorID == null) {
    throw new Error('actorID is null');
  }
  console.log('accepting follow request from', actorID.toString());
  // resolve remote user
  const actor = await fetchPersonByID(actorID);
  const inboxURL = getID(actor.inbox);
  if (inboxURL == null) {
    throw new Error('inboxURL is null');
  }

  // send accept activity
  try {
    const person = buildPerson(c.get('origin'), user);
    const acceptActivity = buildAcceptAcivity(person.id, activity);
    await postActivity(inboxURL, acceptActivity, getPublicKeyID(person.id.toString()), config.privateKey);
  } catch (e) {
    console.error(e);
    throw new Error('Failed to send accept activity');
  }

  // save follower
  try {
    const followersRepo = new UserFollowersRepository();
    const newFollower = RemoteUser.parse(actor);
    await followersRepo.upsert(user, newFollower);
  } catch (e) {
    console.error(e);
    throw new Error('Failed to save follower');
  }
}

export async function deleteFollower(config: Config, c: Context<AppContext>, user: User, activity: UndoActivity) {
  const actorID = getID(activity.actor);
  if (actorID == null) {
    throw new Error('actorID is null');
  }

  console.log('accepting unfollow request from', actorID.toString());
  // resolve remote user
  const actor = await fetchPersonByID(actorID);
  const inboxURL = getID(actor.inbox);
  if (inboxURL == null) {
    throw new Error('inboxURL is null');
  }

  // send accept activity
  try {
    const person = buildPerson(c.get('origin'), user);
    const acceptActivity = buildAcceptAcivity(person.id, activity);
    await postActivity(inboxURL, acceptActivity, getPublicKeyID(person.id.toString()), config.privateKey);
  } catch (e) {
    console.error(e);
    throw new Error('Failed to send accept activity');
  }

  // delete follower
  try {
    const followersRepo = new UserFollowersRepository();
    await followersRepo.delete(user, actorID.toString());
  } catch (e) {
    console.error(e);
    throw new Error('Failed to delete follower');
  }
}
