package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type userRepository struct {
	firestoreClient *firestore.Client
}

func NewUserRepository(firestoreClient *firestore.Client) *userRepository {
	return &userRepository{firestoreClient: firestoreClient}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	userDoc, err := r.firestoreClient.Collection("users").Doc(username).Get(ctx)
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	if err := userDoc.DataTo(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) AddFollower(ctx context.Context, username string, follower *model.Follower) error {
	followers := r.firestoreClient.Collection("users").Doc(username).Collection("followers")
	// check if already exists
	existing := followers.Where("id", "==", follower.ID).Limit(1).Documents(ctx)
	defer existing.Stop()
	if _, err := existing.Next(); err == nil {
		return nil
	}
	// add
	_, _, err := followers.Add(ctx, follower)
	if err != nil {
		return err
	}
	return nil
}
