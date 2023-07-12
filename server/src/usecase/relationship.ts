import { RemoteUser } from '@app/domain/remote-user';
import { User } from '@app/domain/user';
import { UserFollowersRepository } from '@app/repository/user-followers';
import { AppContext } from '@app/web/context';
import { Context } from 'hono';
import { FollowActivity, UndoActivity, buildAcceptAcivity } from '../activitypub/activity';
import { buildPerson, fetchPersonByID } from '../activitypub/person';
import { ActivityPubAgent, getID } from '../activitypub/utilities';

export async function getUserFollowers(user: User): Promise<RemoteUser[]> {
  const followersRepo = new UserFollowersRepository();
  const followers = await followersRepo.list(user);
  return followers;
}

export async function acceptFollowRequest(
  c: Context<AppContext>,
  user: User,
  activity: FollowActivity,
  privateKey: string,
) {
  const actorID = getID(activity.actor);
  if (actorID == null) {
    throw new Error('actorID is null');
  }
  console.log('accept follow request from', actorID);
  // resolve remote user
  const actor = await fetchPersonByID(actorID);

  // send accept activity
  try {
    const agent = new ActivityPubAgent(privateKey);
    const person = buildPerson(c.get('origin'), user);
    const acceptActivity = buildAcceptAcivity(person.id, activity);
    await agent.postActivity(new URL(actorID), person.id.toString(), acceptActivity);
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

export async function deleteFollower(c: Context<AppContext>, user: User, activity: UndoActivity, privateKey: string) {
  const actorID = getID(activity.actor);
  if (actorID == null) {
    throw new Error('actorID is null');
  }

  // send accept activity
  try {
    const agent = new ActivityPubAgent(privateKey);
    const person = buildPerson(c.get('origin'), user);
    const acceptActivity = buildAcceptAcivity(person.id, activity);
    await agent.postActivity(new URL(actorID), person.id.toString(), acceptActivity);
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
