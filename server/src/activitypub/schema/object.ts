import { z } from 'zod';
import { ActivityPubObject, LinkOrURI, URI } from './core';

export const NoteObject = ActivityPubObject.extend({
  type: z.literal('Note'),
  content: z.string(),
  attributedTo: URI,
  to: z.array(URI),
  cc: z.array(URI).optional(),
  inReplyTo: LinkOrURI.optional(),
  sensitive: z.boolean().optional().default(false),
});

export type NoteObject = z.infer<typeof NoteObject>;
