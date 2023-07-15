import { verifyJWT } from '@app/auth/verify';
import { Config } from '@app/domain/config';
import { searchPerson } from '@app/usecase/admin/search-person';
import * as admin from '@app/usecase/admin/users';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>, config: Config) => {
  const adminRoutes = new Hono<AppContext>();

  if (config.isRunningOnCloud) {
    adminRoutes.use('*', verifyJWT());
  }

  adminRoutes.get('/users/list', async (c) => {
    const users = await admin.getUsers();
    return c.json(users);
  });

  adminRoutes.get('/users/:username', async (c) => {
    const username = c.req.param('username');
    const user = await admin.getUserByUsername(username);
    if (user == null) {
      c.status(404);
      return c.json({ error: 'Not Found' });
    }

    return c.json(user);
  });

  adminRoutes.get('/search/person/:resource', handleSearchPerson);

  app.route('/admin', adminRoutes);
};

const handleSearchPerson: Handler<AppContext> = async (c) => {
  const resource = c.req.param('resource');

  try {
    const person = await searchPerson(resource);
    return c.json(person);
  } catch (e) {
    c.status(404);
    return c.json({ error: 'Not Found' });
  }
};
