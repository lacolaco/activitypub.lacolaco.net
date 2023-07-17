import {
  buildCreateActivity,
  buildPerson,
  contextURIsWithExtensions,
  createNoteObject,
  fetchPersonByID,
  getPublicKeyID,
  postActivity,
} from '@app/activitypub';
import { CreateNoteParams, createNote } from '@app/domain/note';
import { User } from '@app/domain/user';
import { UserFollowersRepository } from '@app/repository/user-followers';
import { UserNotesRepository } from '@app/repository/user-notes';
import { getOrigin } from '@app/util/url';
import { KeyObject } from 'node:crypto';

export async function createUserNote(user: User, note: CreateNoteParams, privateKey: KeyObject) {
  console.log('Creating note');
  const noteRepo = new UserNotesRepository();
  const newNote = createNote(note);
  await noteRepo.insert(user, newNote);

  // send Create activity to followers
  const origin = getOrigin(user.host);
  const actor = buildPerson(origin, user);
  if (actor.followers == null) {
    console.log('User', user.id, 'has no followers');
    return;
  }
  const followersRepo = new UserFollowersRepository();
  const followers = await followersRepo.list(user);
  const noteObject = createNoteObject(actor, { ...newNote, cc: [actor.followers] });
  const createActivity = buildCreateActivity(actor.id, noteObject, { contextURIs: contextURIsWithExtensions });
  for (const follower of followers) {
    const inboxPerson = await fetchPersonByID(follower.id);
    if (inboxPerson == null) {
      console.log('Follower', follower.id, 'not found');
      continue;
    }
    console.log('Sending Create activity to', follower.id);
    try {
      await postActivity(inboxPerson.inbox, createActivity, getPublicKeyID(actor.id), privateKey);
    } catch (e) {
      console.log('Error sending Create activity to', follower.id, e);
    }
  }
}
