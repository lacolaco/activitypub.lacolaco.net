package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"google.golang.org/api/iterator"
)

const (
	UsersCollectionName     = "users"
	FollowersCollectionName = "followers"
)

type userRepository struct {
	firestoreClient *firestore.Client
}

func NewUserRepository(firestoreClient *firestore.Client) *userRepository {
	return &userRepository{firestoreClient: firestoreClient}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	userDoc, err := r.firestoreClient.Collection(UsersCollectionName).Doc(username).Get(ctx)
	if err != nil {
		return nil, err
	}
	user, err := model.NewUserFromMap(userDoc.Data())
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Add follorwer to user
func (r *userRepository) AddFollower(ctx context.Context, user *model.User, followerID string) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
	// check if already exists
	existing := col.Where("id", "==", followerID).Limit(1).Documents(ctx)
	defer existing.Stop()
	if _, err := existing.Next(); err == nil {
		return nil
	}
	// add
	data := &model.Follower{ID: followerID, CreatedAt: time.Now()}
	serialized, err := data.ToMap()
	if err != nil {
		return err
	}
	if _, _, err := col.Add(ctx, serialized); err != nil {
		return err
	}
	return nil
}

// Remove follower from user
func (r *userRepository) RemoveFollower(ctx context.Context, user *model.User, followerID string) error {
	followers := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
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

func (r *userRepository) ListFollowers(ctx context.Context, user *model.User) ([]*model.Follower, error) {
	followers := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
	iter := followers.OrderBy("created_at", firestore.Asc).Documents(ctx)
	defer iter.Stop()
	var result []*model.Follower
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		follower, err := model.NewFollowerFromMap(doc.Data())
		if err != nil {
			return nil, err
		}
		result = append(result, follower)
	}
	return result, nil
}
