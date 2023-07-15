import { searchPerson } from '@app/usecase/admin/search-person';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  const adminRoutes = new Hono();

  adminRoutes.get('/users/show/:username', (c) => {
    const headers = Object.fromEntries(c.req.headers.entries());
    console.log(JSON.stringify(headers));

    c.status(404);
    return c.json({ error: 'Not Found' });
  });
  // apiRoutes.get('/search/person/:resource', handleSearchPerson);

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
