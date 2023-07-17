import { randomUUID } from 'node:crypto';
import { z } from 'zod';

export const Note = z.object({
  id: z.string(),
  source: z.string(),
  createdAt: z
    .date()
    .or(z.string())
    .transform((v) => new Date(v)),
  to: z.array(z.string()).optional(),
  cc: z.array(z.string()).optional(),
  inReplyTo: z.string().optional(),
});

export type Note = z.infer<typeof Note>;

export type CreateNoteParams = {
  content: string;
};
export type CreateNoteOptions = {
  id?: string;
  createdAt?: Date;
};

export function createNote(params: CreateNoteParams, options: CreateNoteOptions = {}): Note {
  return Note.parse({
    source: params.content,
    id: options.id ?? randomUUID(),
    createdAt: options.createdAt ?? new Date(),
  });
}
