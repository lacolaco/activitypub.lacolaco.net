import { User } from '@app/domain/user';
import { CollectionReference, Filter, Firestore, Timestamp } from '@google-cloud/firestore';

type UserDocument = {
  id: string;
  host: string;
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

function toDocument(user: User): UserDocument {
  return {
    id: user.id,
    host: user.host,
    username: user.username,
    displayName: user.displayName,
    description: user.description ?? '',
    icon: { url: user.icon.url },
    url: user.url,
    attachments: user.attachments ?? [],
    createdAt: Timestamp.fromDate(user.createdAt),
    updatedAt: Timestamp.fromDate(user.updatedAt),
  };
}

export class UsersRepository {
  readonly db = new Firestore();

  #collection(): CollectionReference<UserDocument> {
    return this.db.collection('users') as CollectionReference<UserDocument>;
  }

  async findByID(uid: string): Promise<User | null> {
    const doc = await this.#collection().doc(uid).get();
    const data = doc.data();
    if (data == null) {
      return null;
    }

    return toUser(data, doc.id);
  }

  async findByUsername(hostname: string, username: string): Promise<User | null> {
    const filter = Filter.and(Filter.where('host', '==', hostname), Filter.where('username', '==', username));
    const query = this.#collection().where(filter).limit(1);
    const items = await query.get();
    if (items.empty) {
      return null;
    }

    const doc = items.docs[0];
    return toUser(doc.data(), doc.id);
  }

  async getUsers(): Promise<User[]> {
    const items = await this.#collection().get();
    return items.docs.map((doc) => toUser(doc.data(), doc.id));
  }

  async insertUser(user: User): Promise<void> {
    const doc = this.#collection().doc(user.id);
    await doc.create(toDocument(user));
  }
}
