import { User } from '@app/domain/user';

export class UsersRepository {
  constructor(readonly db: D1Database) {}

  async findByUsername(username: string): Promise<User | null> {
    const findUser = this.db.prepare(`SELECT * FROM Users WHERE username = ?`).bind(username);
    const user = await findUser.first();
    if (user == null) {
      return null;
    }
    const queryAttachments = this.db.prepare(`SELECT * FROM UserAttachments WHERE userId = ?`).bind(user.id);
    const { results } = await queryAttachments.all();
    user.attachments = results;
    return User.parse(user);
  }
}
