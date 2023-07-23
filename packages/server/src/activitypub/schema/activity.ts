import { z } from 'zod';
import { BaseActivity, TransitiveActivity } from './core';

export const FollowActivity = TransitiveActivity.extend({
  type: z.literal('Follow'),
});
export type FollowActivity = z.infer<typeof FollowActivity>;

export const UndoActivity = TransitiveActivity.extend({
  type: z.literal('Undo'),
});
export type UndoActivity = z.infer<typeof UndoActivity>;

export const AcceptActivity = TransitiveActivity.extend({
  type: z.literal('Accept'),
});
export type AcceptActivity = z.infer<typeof AcceptActivity>;

export const CreateActivity = TransitiveActivity.extend({
  type: z.literal('Create'),
});
export type CreateActivity = z.infer<typeof CreateActivity>;

export const AnyActivity = z.union([FollowActivity, UndoActivity, AcceptActivity, CreateActivity, BaseActivity]);
export type AnyActivity = z.infer<typeof AnyActivity>;
