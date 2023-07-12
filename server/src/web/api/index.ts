import { searchPerson } from '@app/usecase/admin/search-person';
import { Handler, Hono } from 'hono';
import { AppContext } from '../context';

export default (app: Hono<AppContext>) => {
  const apiRoutes = new Hono();
  app.route('/api', apiRoutes);

  // apiRoutes.get('/search/person/:resource', handleSearchPerson);
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
