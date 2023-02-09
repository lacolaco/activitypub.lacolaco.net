package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Client = firestore.Client

func NewFirestoreClient() *firestore.Client {
	c, err := firestore.NewClient(context.Background(), firestore.DetectProjectID)
	if err != nil {
		panic(err)
	}
	return c
}
