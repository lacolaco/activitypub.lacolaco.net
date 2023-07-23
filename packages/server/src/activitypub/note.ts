import { Note } from '@app/domain/note';
import { ActivityPubObject, NoteObject, Person } from './schema';

export function createNoteObject(actor: Person, note: Note, contextURIs?: ActivityPubObject['@context']) {
  const id = `${actor.id}/notes/${note.id}`;
  // TODO: render markdown
  const content = note.source;

  return NoteObject.parse({
    ...(contextURIs && { '@context': contextURIs }),
    id,
    type: 'Note',
    content: content,
    attributedTo: actor.id,
    to: note.to ?? ['https://www.w3.org/ns/activitystreams#Public'],
    cc: note.cc,
    inReplyTo: note.inReplyTo,
    tag: [],
    published: note.createdAt.toISOString(),
  });
}
