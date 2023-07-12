import { RemoteUser } from '@app/domain/remote-user';
import { User } from '@app/domain/user';
import { CollectionReference, Firestore } from '@google-cloud/firestore';

type UserFollowerDocument = {
  id: string;
};

export class UserFollowersRepository {
  readonly db = new Firestore();

  async upsert(user: User, follower: RemoteUser): Promise<void> {
    const collection = this.db
      .collection('users')
      .doc(user.id)
      .collection('followers') as CollectionReference<UserFollowerDocument>;

    const followerDoc = collection.doc(follower.id);
    await followerDoc.set(follower, { merge: true });
  }

  async list(user: User): Promise<RemoteUser[]> {
    const collection = this.db
      .collection('users')
      .doc(user.id)
      .collection('followers') as CollectionReference<UserFollowerDocument>;

    const snapshot = await collection.get();
    const followers = snapshot.docs.map((doc) => {
      const data = doc.data();
      return RemoteUser.parse({
        ...data,
        id: doc.id,
      });
    });

    return followers;
  }
}
