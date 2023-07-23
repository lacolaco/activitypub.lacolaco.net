import { verifyJWT } from '@app/auth/verify';
import { Config } from '@app/domain/config';
import { searchPerson } from '@app/usecase/admin/search-person';
import { createUser, queryUserByUsername, queryUsers } from '@app/usecase/admin/users';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';
import { runInSpan } from '@app/tracing';
import { createUserNote } from '@app/usecase/admin/notes';
import { CreateNoteParams } from '@app/domain/note';
import { NewUserParams } from '@app/domain/user';

export default (app: Hono<AppContext>, config: Config) => {
  const adminRoutes = new Hono<AppContext>();

  if (config.isRunningOnCloud) {
    adminRoutes.use('*', verifyJWT());
  }

  adminRoutes.get('/users', async (c) => {
    return runInSpan('admin.getUsers', async (span) => {
      const users = await queryUsers();
      return c.json(users);
    });
  });

  adminRoutes.post('/users', async (c) => {
    return runInSpan('admin.postUser', async (span) => {
      const body = await c.req.json();
      const parsed = NewUserParams.safeParse(body);
      if (!parsed.success) {
        return c.json({ error: 'Bad Request' }, 400);
      }
      const user = await createUser(parsed.data);
      return c.json(user);
    });
  });

  adminRoutes.get('/users/:hostname/:username', async (c) => {
    return runInSpan('admin.getUser', async (span) => {
      const { hostname, username } = c.req.param();
      span.setAttributes({ hostname, username });
      const user = await queryUserByUsername(hostname, username);
      if (user == null) {
        return c.json({ error: 'Not Found' }, 404);
      }
      return c.json(user);
    });
  });

  adminRoutes.post('/users/:hostname/:username/notes', async (c) => {
    return runInSpan('admin.postUserNotes', async (span) => {
      const { hostname, username } = c.req.param();
      span.setAttributes({ hostname, username });
      const user = await queryUserByUsername(hostname, username);
      if (user == null) {
        return c.json({ error: 'Not Found' }, 404);
      }
      const body = await c.req.json<CreateNoteParams>();
      await createUserNote(user, body, c.get('config').privateKey);
      return c.json(user);
    });
  });

  adminRoutes.get('/search/person/:resource', handleSearchPerson);

  app.route('/admin', adminRoutes);
};

const handleSearchPerson: Handler<AppContext> = async (c) => {
  return runInSpan('admin.searchPerson', async (span) => {
    const { resource } = c.req.param();
    span.setAttributes({ resource });
    try {
      const person = await searchPerson(resource);
      return c.json(person);
    } catch (e) {
      return c.json({ error: 'Not Found' }, 404);
    }
  });
};
