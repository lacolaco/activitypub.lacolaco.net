import { z } from 'zod';

export const RemoteUser = z.object({
  id: z.string(),
});

export type RemoteUser = z.infer<typeof RemoteUser>;
