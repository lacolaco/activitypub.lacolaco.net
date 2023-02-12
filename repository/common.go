package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Item interface {
	GetID() string
}

// add item to collection if not exists
func addIfNotExists(ctx context.Context, collection *firestore.CollectionRef, item Item) error {
	// check if already exists
	existing := collection.Where("id", "==", item.GetID()).Limit(1).Documents(ctx)
	defer existing.Stop()
	if _, err := existing.Next(); err == nil {
		return nil
	}
	// add
	if _, _, err := collection.Add(ctx, item); err != nil {
		return err
	}
	return nil
}

// remove item by id from collection
func removeIfExists(ctx context.Context, collection *firestore.CollectionRef, id string) error {
	// check if already exists
	existing := collection.Where("id", "==", id).Limit(1).Documents(ctx)
	defer existing.Stop()
	doc, err := existing.Next()
	if err != nil {
		return nil
	}
	// remove
	if _, err := doc.Ref.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// get items from collection
func getAllItems[T interface{}](ctx context.Context, collection *firestore.CollectionRef, q firestore.Query) ([]T, error) {
	iter := collection.Documents(ctx)
	defer iter.Stop()
	var result []T
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item T
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}
