import { ActivityStreamsObject, ObjectOrLinkOrURI, OrderedCollection, URI } from './schema';

export function buildOrderedCollection(
  collectionID: URI,
  items: ObjectOrLinkOrURI[],
  contextURIs: ActivityStreamsObject['@context'] = [],
) {
  return OrderedCollection.parse({
    ...(contextURIs && { '@context': contextURIs }),
    id: collectionID,
    type: 'OrderedCollection',
    totalItems: items.length,
    orderedItems: items,
  });
}
