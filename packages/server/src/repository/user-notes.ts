import { Note } from '@app/domain/note';
import { User } from '@app/domain/user';
import { CollectionReference, Firestore, Timestamp } from '@google-cloud/firestore';

type UserNoteDocument = {
  id: string;
  createdAt: Timestamp;
};

function toNote(doc: UserNoteDocument): Note {
  return Note.parse({
    ...doc,
    createdAt: doc.createdAt.toDate(),
  });
}

export class UserNotesRepository {
  readonly #db = new Firestore();

  #collection(user: User): CollectionReference<UserNoteDocument> {
    return this.#db.collection('users').doc(user.id).collection('notes') as CollectionReference<UserNoteDocument>;
  }

  async insert(user: User, note: Note): Promise<void> {
    const collection = this.#collection(user);
    const query = collection.where('id', '==', note.id).limit(1);
    const snapshot = await query.get();
    if (!snapshot.empty) {
      throw new Error(`Note ${note.id} already exists`);
    }
    await collection.add({
      ...note,
      createdAt: Timestamp.fromDate(note.createdAt),
    });
  }

  async list(user: User, limit = 20): Promise<Note[]> {
    const query = this.#collection(user).orderBy('createdAt', 'desc').limit(limit);

    const snapshot = await query.get();
    const notes = snapshot.docs.map((doc) => {
      const data = doc.data();
      return toNote(data);
    });

    return notes;
  }
}
