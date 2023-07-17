import { verifyJWT } from '@app/auth/verify';
import { Config } from '@app/domain/config';
import { searchPerson } from '@app/usecase/admin/search-person';
import * as admin from '@app/usecase/admin/users';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';
import { runInSpan } from '@app/tracing';

export default (app: Hono<AppContext>, config: Config) => {
  const adminRoutes = new Hono<AppContext>();

  if (config.isRunningOnCloud) {
    adminRoutes.use('*', verifyJWT());
  }

  adminRoutes.get('/users', async (c) => {
    return runInSpan('admin.getUsers', async (span) => {
      const users = await admin.getUsers();
      return c.json(users);
    });
  });

  adminRoutes.get('/users/:hostname/:username', async (c) => {
    return runInSpan('admin.getUser', async (span) => {
      const { hostname, username } = c.req.param();
      span.setAttributes({ hostname, username });
      const user = await admin.getUserByUsername(hostname, username);
      if (user == null) {
        return c.json({ error: 'Not Found' }, 404);
      }
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
