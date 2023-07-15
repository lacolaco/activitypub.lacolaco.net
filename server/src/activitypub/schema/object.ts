import { z } from 'zod';
import { ActivityStreamsObject } from './core';

export const Note = ActivityStreamsObject.extend({
  type: z.literal('Note'),
});

export type Note = z.infer<typeof Note>;
