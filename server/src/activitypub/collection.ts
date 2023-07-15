import { contextURIs } from './context';
import { ObjectOrLinkOrURI, OrderedCollection, URI } from './schema';

export function buildOrderedCollection(collectionID: URI, items: ObjectOrLinkOrURI[]) {
  return OrderedCollection.parse({
    '@context': contextURIs,
    id: collectionID,
    type: 'OrderedCollection',
    totalItems: items.length,
    orderedItems: items,
  });
}
