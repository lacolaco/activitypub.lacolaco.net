import { UsersRepository } from '@app/repository/users';
import { getTracer } from '@app/tracing';

export async function getUsers() {
  return getTracer().startActiveSpan('admin.getUsers', async (span) => {
    const userRepo = new UsersRepository();
    const users = await userRepo.getUsers();
    return users;
  });
}

export async function getUserByUsername(hostname: string, username: string) {
  return getTracer().startActiveSpan('admin.getUserByUsername', async (span) => {
    const userRepo = new UsersRepository();
    const user = await userRepo.findByUsername(hostname, username);
    if (user == null) {
      return null;
    }
    return user;
  });
}
