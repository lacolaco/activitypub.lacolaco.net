import { describe, expect, test } from 'vitest';
import { AcceptActivity, AnyActivity, FollowActivity, UndoActivity } from './activity';
import { ActivityPubObject, ObjectOrLinkOrURI, Person } from './core';

describe('activitypub schema', () => {
  describe('Activity', () => {
    test('mastodon-follow as FollowActivity', async () => {
      const json = await import('./fixtures/mastodon-follow.json');
      const parsed = FollowActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('mastodon-undo-follow as UndoActivity', async () => {
      const json = await import('./fixtures/mastodon-undo-follow.json');
      const parsed = UndoActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('misskey-accept-follow as AcceptActivity', async () => {
      const json = await import('./fixtures/misskey-accept-follow.json');
      const parsed = AcceptActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('misskey-accept-follow as AnyActivity', async () => {
      const json = await import('./fixtures/misskey-accept-follow.json');
      const parsed = AnyActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('misskey-delete-person as AnyActivity', async () => {
      const json = await import('./fixtures/misskey-delete-person.json');
      const parsed = AnyActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('mastodon-delete-person as AnyActivity', async () => {
      const json = await import('./fixtures/mastodon-delete-person.json');
      const parsed = AnyActivity.parse(json);
      expect(parsed).toBeTruthy();
    });

    test('misskey-create-note as AnyActivity', async () => {
      const json = await import('./fixtures/misskey-create-note.json');
      const parsed = AnyActivity.parse(json);
      expect(parsed).toBeTruthy();
    });
  });

  describe('Person', () => {
    test('misskey-person as Object', async () => {
      const json = await import('./fixtures/misskey-person.json');
      const parsed = ActivityPubObject.parse(json);
      expect(parsed).toBeTruthy();
    });
    test('misskey-person as ObjectOrLinkOrURI', async () => {
      const json = await import('./fixtures/misskey-person.json');
      const parsed = ObjectOrLinkOrURI.parse(json);
      expect(parsed).toBeTruthy();
    });
    test('misskey-person as Person', async () => {
      const json = await import('./fixtures/misskey-person.json');
      const parsed = Person.parse(json);
      expect(parsed).toBeTruthy();
    });
  });
});
