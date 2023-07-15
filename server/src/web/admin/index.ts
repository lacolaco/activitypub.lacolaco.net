import { verifyJWT } from '@app/auth/verify';
import { Config } from '@app/domain/config';
import { searchPerson } from '@app/usecase/admin/search-person';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';
import { UsersRepository } from '@app/repository/users';

export default (app: Hono<AppContext>, config: Config) => {
  const adminRoutes = new Hono<AppContext>();

  if (config.isRunningOnCloud) {
    adminRoutes.use('*', verifyJWT());
  }

  adminRoutes.get('/users/show/:username', async (c) => {
    const username = c.req.param('username');
    const userRepo = new UsersRepository();
    const user = await userRepo.findByUsername(username);
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
