import { RemoteUser } from '@app/domain/remote-user';
import { User } from '@app/domain/user';
import { CollectionReference, Firestore } from '@google-cloud/firestore';

type UserFollowerDocument = {
  id: string;
  createdAt: Date;
};

export class UserFollowersRepository {
  readonly db = new Firestore();

  async upsert(user: User, follower: RemoteUser): Promise<void> {
    const collection = this.db
      .collection('users')
      .doc(user.id)
      .collection('followers') as CollectionReference<UserFollowerDocument>;

    const query = collection.where('id', '==', follower.id).limit(1);
    const snapshot = await query.get();
    if (snapshot.empty) {
      await collection.add({
        id: follower.id,
        createdAt: new Date(),
      });
    } else {
      const doc = snapshot.docs[0];
      await doc.ref.update(follower);
    }
  }

  async list(user: User): Promise<RemoteUser[]> {
    const collection = this.db
      .collection('users')
      .doc(user.id)
      .collection('followers') as CollectionReference<UserFollowerDocument>;

    const query = collection.orderBy('createdAt', 'desc');

    const snapshot = await query.get();
    const followers = snapshot.docs.map((doc) => {
      const data = doc.data();
      return RemoteUser.parse({
        ...data,
      });
    });

    return followers;
  }
}
