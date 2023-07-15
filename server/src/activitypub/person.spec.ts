import { createUser } from '@app/domain/testing/user';
import { describe, expect, test } from 'vitest';
import { buildPerson } from './person';

describe('buildPerson', () => {
  test('build a person', async () => {
    const user = createUser();

    const person = buildPerson('http://localhost', user);

    expect(JSON.parse(JSON.stringify(person))).toMatchInlineSnapshot(`
      {
        "@context": [
          "https://www.w3.org/ns/activitystreams",
          "https://w3id.org/security/v1",
          {
            "Emoji": "toot:Emoji",
            "Hashtag": "as:Hashtag",
            "PropertyValue": "schema:PropertyValue",
            "discoverable": "toot:discoverable",
            "featured": "toot:featured",
            "manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
            "misskey": "https://misskey-hub.net/ns#",
            "quoteUrl": "as:quoteUrl",
            "schema": "http://schema.org#",
            "sensitive": "as:sensitive",
            "toot": "http://joinmastodon.org/ns#",
            "value": "schema:value",
          },
        ],
        "attachment": [],
        "discoverable": true,
        "endpoints": {
          "sharedInbox": "http://localhost/inbox",
        },
        "followers": "http://localhost/users/1/followers",
        "following": "http://localhost/users/1/following",
        "icon": {
          "type": "Image",
          "url": "https://example.com/avatar.png",
        },
        "id": "http://localhost/users/1",
        "inbox": "http://localhost/users/1/inbox",
        "manuallyApprovesFollowers": false,
        "name": "test",
        "outbox": "http://localhost/users/1/outbox",
        "preferredUsername": "test",
        "published": "2006-01-02T15:04:05.999Z",
        "summary": "test",
        "type": "Person",
        "updated": "2006-01-02T15:04:05.999Z",
        "url": "https://example.com/@test",
      }
    `);
  });
});
