import { z } from 'zod';

export const RemoteUser = z
  .object({
    id: z.string(),
  })
  .passthrough();

export type RemoteUser = z.infer<typeof RemoteUser>;
