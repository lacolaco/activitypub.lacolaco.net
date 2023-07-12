import { RemoteUser } from '@app/domain/remote-user';
import { User } from '@app/domain/user';
import { UserFollowersRepository } from '@app/repository/user-followers';
import { FollowActivity, UndoActivity, buildAcceptAcivity, getActorOf } from '../activitypub/activity';
import { fetchPersonByID } from '../activitypub/person';
import { ActivityPubAgent, getEntityID } from '../activitypub/utilities';

export async function getUserFollowers(user: User): Promise<RemoteUser[]> {
  const followersRepo = new UserFollowersRepository();
  const followers = await followersRepo.list(user);
  return followers;
}

export async function acceptFollowRequest(user: User, activity: FollowActivity, privateKey: string) {
  const actorID = getEntityID(getActorOf(activity));
  if (actorID == null) {
    throw new Error('actorID is null');
  }
  // resolve remote user
  const actor = await fetchPersonByID(actorID);

  // send accept activity
  try {
    const agent = new ActivityPubAgent(privateKey);
    const acceptActivity = buildAcceptAcivity(new URL(user.id), activity);
    await agent.postActivity(new URL(actorID), user.id, acceptActivity);
  } catch (e) {
    throw new Error('Failed to send accept activity');
  }

  // save follower
  try {
    const followersRepo = new UserFollowersRepository();
    const newFollower = RemoteUser.parse(actor);
    await followersRepo.upsert(user, newFollower);
  } catch (e) {
    throw new Error('Failed to save follower');
  }
}

export async function deleteFollower(user: User, activity: UndoActivity, privateKey: string) {
  const actorID = getEntityID(getActorOf(activity));
  if (actorID == null) {
    throw new Error('actorID is null');
  }

  // send accept activity
  try {
    const agent = new ActivityPubAgent(privateKey);
    const acceptActivity = buildAcceptAcivity(new URL(user.id), activity);
    await agent.postActivity(new URL(actorID), user.id, acceptActivity);
  } catch (e) {
    throw new Error('Failed to send accept activity');
  }

  // delete follower
  try {
    const followersRepo = new UserFollowersRepository();
    await followersRepo.delete(user, actorID.toString());
  } catch (e) {
    throw new Error('Failed to delete follower');
  }
}
