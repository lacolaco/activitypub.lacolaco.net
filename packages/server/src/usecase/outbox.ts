import { User } from '@app/domain/user';
import { UserNotesRepository } from '@app/repository/user-notes';
import {
  OrderedCollection,
  buildOrderedCollection,
  buildPerson,
  contextURIsWithExtensions,
  createNoteObject,
} from '@app/activitypub';

export async function buildOutbox(origin: string, user: User): Promise<OrderedCollection> {
  const notesRepo = new UserNotesRepository();
  const notes = await notesRepo.list(user);
  const actor = buildPerson(origin, user);
  const noteObjects = notes.map((note) => createNoteObject(actor, note));
  return buildOrderedCollection(actor.outbox, noteObjects, contextURIsWithExtensions);
}
