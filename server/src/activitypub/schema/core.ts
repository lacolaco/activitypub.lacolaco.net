import { z } from 'zod';

export const URI = z
  .string()
  .or(z.instanceof(URL))
  .transform((v) => v.toString());
export type URI = z.infer<typeof URI>;

const JSONLDObject = z.object({
  '@context': z.any().or(z.array(z.any())).optional(),
});

// https://www.w3.org/TR/activitystreams-core/#object
export const ActivityStreamsObject = JSONLDObject.extend({
  '@context': z.any().or(z.array(z.any())).optional(),
  type: z.string().optional(),
  id: URI.optional(),
}).passthrough();
export type ActivityStreamsObject = z.infer<typeof ActivityStreamsObject>;

export const ActivityStreamsLink = JSONLDObject.extend({
  href: URI,
  name: z.string(),
});
export type ActivityStreamsLink = z.infer<typeof ActivityStreamsLink>;

export const LinkOrURI = z.union([ActivityStreamsLink, URI]);
export type LinkOrURI = z.infer<typeof LinkOrURI>;
export const ObjectOrURI = z.union([ActivityStreamsObject, URI]);
export type ObjectOrURI = z.infer<typeof ObjectOrURI>;
export const ObjectOrLink = z.union([ActivityStreamsObject, ActivityStreamsLink]);
export type ObjectOrLink = z.infer<typeof ObjectOrLink>;
export const ObjectOrLinkOrURI = z.union([ActivityStreamsLink, ActivityStreamsObject, URI]);
export type ObjectOrLinkOrURI = z.infer<typeof ObjectOrLinkOrURI>;

export const ObjectRef = z.union([URI, ActivityStreamsObject]);
export type ObjectRef = z.infer<typeof ObjectRef>;

export const BaseActivity = ActivityStreamsObject.extend({
  actor: ObjectRef,
  object: ObjectRef.optional(),
  target: ObjectRef.optional(),
});
export type BaseActivity = z.infer<typeof BaseActivity>;

export const TransitiveActivity = BaseActivity.extend({
  object: ObjectRef,
});
export type TransitiveActivity = z.infer<typeof TransitiveActivity>;

export const IntransitiveActivity = BaseActivity.extend({});
export type IntransitiveActivity = z.infer<typeof IntransitiveActivity>;

export const ImageObject = ActivityStreamsObject.extend({
  type: z.literal('Image'),
  url: LinkOrURI.or(z.array(LinkOrURI)),
});

const DateTime = z
  .date()
  .or(z.string())
  .transform((v) => new Date(v));

// https://www.w3.org/TR/activitypub/#obj
export const ActivityPubObject = ActivityStreamsObject.extend({
  // MUST have
  id: URI,
  type: z.string(),
  // MAY have
  icon: ImageObject.optional(),
  name: z.string().optional(),
  summary: z.string().optional(),
  url: URI.optional(),
  published: DateTime.optional(),
  updated: DateTime.optional(),
  attachment: z.array(ActivityStreamsObject).optional(),
});

/**
 * Public key for HTTP Signatures.
 *
 * @see https://docs.joinmastodon.org/spec/activitypub/#publicKey
 */
export const PublicKey = z.object({
  id: URI,
  owner: URI,
  publicKeyPem: z.string(),
});
export type PublicKey = z.infer<typeof PublicKey>;

/**
 * https://www.w3.org/TR/activitypub/#actor-objects
 */
export const Person = ActivityPubObject.extend({
  // MUST have
  type: z.literal('Person'),
  inbox: URI,
  outbox: URI,
  // SHOULD have
  following: URI.optional(),
  followers: URI.optional(),
  // MAY have
  liked: URI.optional(),
  preferredUsername: z.string().optional(),
  endpoints: z
    .object({
      sharedInbox: URI.optional(),
    })
    .optional(),
  // extensions
  alsoKnownAs: z.array(URI).optional(),
  discoverable: z.boolean().optional(),
  featured: URI.optional(),
  manuallyApprovesFollowers: z.boolean().optional(),
  publicKey: PublicKey.optional(),
});
export type Person = z.infer<typeof Person>;

export const Collection = ActivityPubObject.extend({
  type: z.literal('Collection'),
  totalItems: z.number(),
  items: z.array(ObjectOrLinkOrURI),
});

export const OrderedCollection = Collection.extend({
  type: z.literal('OrderedCollection'),
  orderedItems: z.array(ObjectOrLinkOrURI),
});
export type OrderedCollection = z.infer<typeof OrderedCollection>;
