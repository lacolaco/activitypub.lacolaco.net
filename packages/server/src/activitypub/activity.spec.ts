import { describe, expect, test } from 'vitest';
import { buildAcceptAcivity } from './activity';
import { FollowActivity, URI } from './schema';

describe('buildAcceptAcivity', () => {
  test('returns an Accept activity', () => {
    const actorID = URI.parse('https://example.com/users/alice');
    const object = FollowActivity.parse({
      '@context': 'https://www.w3.org/ns/activitystreams',
      type: 'Follow',
      id: URI.parse('https://example.com/users/alice/follows/123'),
      actor: URI.parse('https://example.com/users/alice'),
      object: URI.parse('https://example.com/users/bob'),
    });
    const accept = buildAcceptAcivity(actorID, object);

    expect(accept.type).toStrictEqual('Accept');
    expect(accept.id?.toString().startsWith('https://example.com/users/alice/accept/')).toBe(true);
    expect(accept.actor).toStrictEqual(actorID);
    expect(accept.object).toStrictEqual(object);
  });
});
