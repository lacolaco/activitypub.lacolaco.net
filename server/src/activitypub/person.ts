import { AP, isTypeOf } from '@activity-kit/types';
import { User } from '@app/domain/user';
import { contextURIs } from './context';
import { buildPropertyValue } from './property-value';
import { getPublicKeyID } from './signature';

type Person = AP.Person & Record<string, unknown>;

export function buildPerson(origin: string, user: User, publicKey: string) {
  // `preferredUsername` はあとから変更可能にするため、不変なURIを `id` として使う
  const userURI = `${origin}/users/${user.id}`;

  return {
    '@context': contextURIs,
    id: new URL(userURI),
    type: 'Person',
    name: user.displayName,
    preferredUsername: user.username,
    summary: user.description ?? '',
    icon: {
      type: 'Image',
      url: new URL(user.icon.url),
    },
    attachment: user.attachments?.map((a) => buildPropertyValue(a.name, a.value)) as unknown as AP.EntityReference[],
    inbox: new URL(`${userURI}/inbox`),
    outbox: new URL(`${userURI}/outbox`),
    followers: new URL(`${userURI}/followers`),
    following: new URL(`${userURI}/following`),
    endpoints: {
      sharedInbox: new URL(`${origin}/inbox`),
    },
    url: new URL(userURI), // TODO: use user's website
    published: user.createdAt,
    updated: user.updatedAt,
    manuallyApprovesFollowers: false,
    discoverable: true,
    publicKey: {
      id: getPublicKeyID(userURI),
      owner: userURI,
      publicKeyPem: publicKey,
    },
  } satisfies Person;
}

export async function fetchPersonByID(id: URL): Promise<Person> {
  const res = await fetch(id, {
    method: 'GET',
    headers: {
      Accept: 'application/activity+json, application/json',
    },
  });
  if (res.status !== 200) {
    throw new Error(`unexpected status code: ${res.status}`);
  }

  const person = await res.json();
  if (!isTypeOf<Person>(person, AP.ActorTypes)) {
    throw new Error(`unexpected type: ${person.type}`);
  }
  return person;
}
