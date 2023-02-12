package repository

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var (
	ErrNotFound = errors.New("not found")
)

type Item interface {
	GetID() string
}

func deleteItems(ctx context.Context, q firestore.Query) error {
	iter := q.Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err != nil {
		return nil
	}
	if _, err := doc.Ref.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func findItem[T interface{}](ctx context.Context, q firestore.Query) (*T, error) {
	iter := q.Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var item T
	if err := doc.DataTo(&item); err != nil {
		return nil, err
	}
	return &item, nil

}

func getAllItems[T interface{}](ctx context.Context, q firestore.Query) ([]*T, error) {
	iter := q.Documents(ctx)
	defer iter.Stop()
	var result []*T
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
		result = append(result, &item)
	}
	return result, nil
}
