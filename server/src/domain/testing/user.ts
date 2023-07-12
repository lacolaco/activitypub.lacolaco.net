import { User } from '../user';

export function createUser(override?: Partial<User>): User {
  return {
    id: '1',
    username: 'test',
    displayName: 'test',
    description: 'test',
    icon: {
      url: 'https://example.com/avatar.png',
    },
    url: 'https://example.com/banner.png',
    createdAt: new Date('2006-01-02T15:04:05.999Z'),
    updatedAt: new Date('2006-01-02T15:04:05.999Z'),
    attachments: [],
    ...override,
  };
}
