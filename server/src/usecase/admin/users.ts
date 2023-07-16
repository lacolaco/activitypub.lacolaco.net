import { UsersRepository } from '@app/repository/users';
import { runInSpan } from '@app/tracing';

export async function getUsers() {
  return runInSpan('admin.getUsers', async (span) => {
    try {
      const userRepo = new UsersRepository();
      const users = await userRepo.getUsers();
      return users;
    } finally {
      span.end();
    }
  });
}

export async function getUserByUsername(hostname: string, username: string) {
  return runInSpan('admin.getUserByUsername', async (span) => {
    try {
      const userRepo = new UsersRepository();
      const user = await userRepo.findByUsername(hostname, username);
      if (user == null) {
        return null;
      }
      return user;
    } finally {
      span.end();
    }
  });
}
