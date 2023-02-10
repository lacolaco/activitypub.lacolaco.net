package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"google.golang.org/api/iterator"
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

// Add follorwer to user
func (r *userRepository) AddFollower(ctx context.Context, username, followerID string) error {
	followers := r.firestoreClient.Collection("users").Doc(username).Collection("followers")
	// check if already exists
	existing := followers.Where("id", "==", followerID).Limit(1).Documents(ctx)
	defer existing.Stop()
	if _, err := existing.Next(); err == nil {
		return nil
	}
	// add
	_, _, err := followers.Add(ctx, map[string]interface{}{
		"id":         followerID,
		"created_at": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

// Remove follower from user
func (r *userRepository) RemoveFollower(ctx context.Context, username, followerID string) error {
	followers := r.firestoreClient.Collection("users").Doc(username).Collection("followers")
	// check if already exists
	existing := followers.Where("id", "==", followerID).Limit(1).Documents(ctx)
	defer existing.Stop()
	doc, err := existing.Next()
	if err == iterator.Done {
		// not found
		return nil
	}
	if err != nil {
		return err
	}
	// remove
	if _, err := doc.Ref.Delete(ctx); err != nil {
		return err
	}
	return nil
}
