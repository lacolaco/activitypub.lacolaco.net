import { User } from '@app/domain/user';
import { CollectionReference, Filter, Firestore, Timestamp } from '@google-cloud/firestore';

type UserDocument = {
  id: string;
  username: string;
  displayName: string;
  description: string;
  icon: { url: string };
  url: string;
  attachments: Array<{ name: string; value: string }>;
  createdAt: Timestamp;
  updatedAt: Timestamp;
};

function toUser(doc: UserDocument, id: string): User {
  return User.parse({
    ...doc,
    id,
    createdAt: doc.createdAt.toDate(),
    updatedAt: doc.updatedAt.toDate(),
  });
}

export class UsersRepository {
  readonly db = new Firestore();

  async findByID(uid: string): Promise<User | null> {
    const collection = this.db.collection('users') as CollectionReference<UserDocument>;
    const doc = await collection.doc(uid).get();
    const data = doc.data();
    if (data == null) {
      return null;
    }

    return toUser(data, doc.id);
  }

  async findByUsername(hostname: string, username: string): Promise<User | null> {
    const collection = this.db.collection('users') as CollectionReference<UserDocument>;
    const filter = Filter.and(Filter.where('host', '==', hostname), Filter.where('username', '==', username));
    const query = collection.where(filter).limit(1);
    const items = await query.get();
    if (items.empty) {
      return null;
    }

    const doc = items.docs[0];
    return toUser(doc.data(), doc.id);
  }

  async getUsers(): Promise<User[]> {
    const collection = this.db.collection('users') as CollectionReference<UserDocument>;
    const items = await collection.get();
    return items.docs.map((doc) => toUser(doc.data(), doc.id));
  }
}
