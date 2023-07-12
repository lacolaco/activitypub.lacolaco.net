import { User } from '@app/domain/user';
import { CollectionReference, Firestore, Timestamp } from '@google-cloud/firestore';

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

export class UsersRepository {
  readonly db = new Firestore();

  async findByID(uid: string): Promise<User | null> {
    const collection = this.db.collection('users') as CollectionReference<UserDocument>;
    const userDoc = await collection.doc(uid).get();
    const data = userDoc.data();
    if (data == null) {
      return null;
    }

    return User.parse({
      ...data,
      id: userDoc.id,
      createdAt: data.createdAt.toDate(),
      updatedAt: data.updatedAt.toDate(),
    });
  }

  async findByUsername(username: string): Promise<User | null> {
    const collection = this.db.collection('users') as CollectionReference<UserDocument>;
    const query = collection.where('username', '==', username).limit(1);
    const snapshot = await query.get();
    if (snapshot.empty) {
      return null;
    }

    const userDoc = snapshot.docs[0];
    const data = userDoc.data();
    return User.parse({
      ...data,
      id: userDoc.id,
      createdAt: data.createdAt.toDate(),
      updatedAt: data.updatedAt.toDate(),
    });
  }
}
