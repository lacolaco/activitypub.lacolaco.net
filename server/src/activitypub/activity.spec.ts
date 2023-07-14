import { AP } from '@activity-kit/types';
import { describe, expect, test } from 'vitest';
import { buildAcceptAcivity } from './activity';

describe('buildAcceptAcivity', () => {
  test('returns an Accept activity', () => {
    const actorID = new URL('https://example.com/users/alice');
    const object = {
      '@context': 'https://www.w3.org/ns/activitystreams',
      type: 'Follow',
      id: new URL('https://example.com/users/alice/follows/123'),
      actor: new URL('https://example.com/users/alice'),
      object: new URL('https://example.com/users/bob'),
    } satisfies AP.Follow;
    const accept = buildAcceptAcivity(actorID, object);

    expect(accept.id.toString().startsWith('https://example.com/users/alice/accept/')).toBe(true);
    expect(accept).toMatchObject({
      type: 'Accept',
      actor: actorID,
      object: object,
    });
  });
});
