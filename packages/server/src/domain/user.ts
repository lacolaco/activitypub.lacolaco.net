import { randomUUID } from 'crypto';
import { z } from 'zod';

export const User = z.object({
  id: z.string(),
  host: z.string(),
  username: z.string(),
  displayName: z.string(),
  description: z.string().nullable(),
  icon: z.object({
    url: z.string(),
  }),
  url: z.string(),
  attachments: z
    .array(
      z.object({
        name: z.string(),
        value: z.string(),
      }),
    )
    .optional(),
  createdAt: z
    .date()
    .or(z.string())
    .transform((v) => new Date(v)),
  updatedAt: z
    .date()
    .or(z.string())
    .transform((v) => new Date(v)),
});

export type User = z.infer<typeof User>;

export const NewUserParams = z.object({
  id: z.string().optional(),
  host: z.string(),
  username: z.string(),
  displayName: z.string(),
  description: z.string().optional(),
  icon: z.object({
    url: z.string(),
  }),
  url: z.string(),
  attachments: z
    .array(
      z.object({
        name: z.string(),
        value: z.string(),
      }),
    )
    .optional(),
  createdAt: z.date().optional(),
});
export type NewUserParams = z.infer<typeof NewUserParams>;

export function newUser(params: NewUserParams): User {
  const now = params.createdAt ?? new Date();
  return User.parse({
    ...params,
    id: params.id ?? randomUUID(),
    createdAt: now,
    updatedAt: now,
  });
}
