import { AP } from '@activity-kit/types';
import { contextURIs } from './context';

export function buildOrderedCollection(id: URL, items: AP.EntityReference[]): AP.OrderedCollection {
	return {
		'@context': contextURIs,
		id,
		type: 'OrderedCollection',
		totalItems: items.length,
		orderedItems: items,
	};
}
