import { z } from 'zod';

export const User = z.object({
	id: z.number(),
	username: z.string(),
	displayName: z.string(),
	description: z.string().nullable(),
	iconUrl: z.string(),
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
