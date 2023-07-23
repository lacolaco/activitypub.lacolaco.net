import { User } from '@app/domain/user';
import { contextURIsWithExtensions } from './context';
import { buildPropertyValue } from './property-value';
import { Person, URI } from './schema';
import { getPublicKeyID } from './signature';

export function buildPerson(origin: string, user: User, publicKey?: string) {
  // `preferredUsername` はあとから変更可能にするため、不変なURIを `id` として使う
  const userURI = `${origin}/users/${user.id}`;

  return Person.parse({
    '@context': contextURIsWithExtensions,
    id: userURI,
    type: 'Person',
    name: user.displayName,
    preferredUsername: user.username,
    summary: user.description ?? '',
    icon: {
      type: 'Image',
      url: user.icon.url,
    },
    inbox: `${userURI}/inbox`,
    outbox: `${userURI}/outbox`,
    followers: `${userURI}/followers`,
    following: `${userURI}/following`,
    endpoints: {
      sharedInbox: new URL(`${origin}/inbox`),
    },
    url: user.url,
    published: user.createdAt,
    updated: user.updatedAt,
    attachment: user.attachments?.map((a) => buildPropertyValue(a.name, a.value)),
    manuallyApprovesFollowers: false,
    discoverable: true,
    publicKey: publicKey && {
      id: getPublicKeyID(userURI),
      owner: userURI,
      publicKeyPem: publicKey,
    },
  });
}

export async function fetchPersonByID(id: URI): Promise<Person> {
  const res = await fetch(id, {
    method: 'GET',
    headers: { Accept: 'application/activity+json, application/json' },
  });
  if (res.status !== 200) {
    throw new Error(`unexpected status code: ${res.status}`);
  }
  const body = await res.json();
  const parsed = Person.safeParse(body);
  if (parsed.success === false) {
    console.error(parsed.error);
    throw new Error(`unexpected object`);
  }
  return parsed.data;
}
