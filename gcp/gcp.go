package gcp

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
)

func NewFirestoreClient() *firestore.Client {
	c, err := firestore.NewClient(context.Background(), firestore.DetectProjectID)
	if err != nil {
		panic(err)
	}
	return c
}

func NewFirebaseAuthClient() *firebaseauth.Client {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	c, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}
	return c
}
