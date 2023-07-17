import { CreateNoteParams, createNote } from '@app/domain/note';
import { User } from '@app/domain/user';
import { UserNotesRepository } from '@app/repository/user-notes';

export async function createUserNote(user: User, note: CreateNoteParams) {
  console.log('Creating note');
  const noteRepo = new UserNotesRepository();
  const newNote = createNote(note);
  await noteRepo.insert(user, newNote);
}
