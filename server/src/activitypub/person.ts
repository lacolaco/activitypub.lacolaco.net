import { AP } from '@activity-kit/types';
import { User } from '@app/domain/user';
import { contextURIs } from './context';
import { buildPropertyValue } from './property-value';

type Person = AP.Person & Record<string, unknown>;

export function buildPerson(origin: string, user: User) {
  const userURI = `${origin}/ap/users/${user.username}`;
  return {
    '@context': contextURIs,
    id: new URL(userURI),
    type: 'Person',
    name: user.displayName,
    preferredUsername: user.username,
    summary: user.description ?? '',
    icon: {
      type: 'Image',
      url: new URL(user.iconUrl),
    },
    attachment: user.attachments?.map((a) => buildPropertyValue(a.name, a.value)) as unknown as AP.EntityReference[],
    inbox: new URL(`${userURI}/inbox`),
    outbox: new URL(`${userURI}/outbox`),
    followers: new URL(`${userURI}/followers`),
    following: new URL(`${userURI}/following`),
    endpoints: {
      sharedInbox: new URL(`${origin}/ap/inbox`),
    },
    url: new URL(userURI), // TODO: use user's website
    published: user.createdAt,
    updated: user.updatedAt,
    manuallyApprovesFollowers: false,
    discoverable: true,
  } satisfies Person;
}
