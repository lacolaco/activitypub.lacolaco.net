package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

const (
	UsersCollectionName     = "users"
	FollowersCollectionName = "followers"
	FollowingCollectionName = "following"
)

type userRepo struct {
	firestoreClient *firestore.Client
}

func NewUserRepository(firestoreClient *firestore.Client) *userRepo {
	return &userRepo{firestoreClient: firestoreClient}
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.LocalUser, error) {
	collection := r.firestoreClient.Collection(UsersCollectionName)
	users, err := getAllItems[*model.LocalUser](ctx, collection, collection.Where("preferred_username", "==", username))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found: %s", username)
	}
	return users[0], nil
}

func (r *userRepo) FindByID(ctx context.Context, uid string) (*model.LocalUser, error) {
	userDoc, err := r.firestoreClient.Collection(UsersCollectionName).Doc(uid).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user *model.LocalUser
	if err := userDoc.DataTo(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// Add follorwer to user
func (r *userRepo) AddFollower(ctx context.Context, user *model.LocalUser, followerID string) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
	data := &model.RemoteUser{ID: followerID, CreatedAt: time.Now()}
	if err := addIfNotExists(ctx, col, data); err != nil {
		return err
	}
	return nil
}

// Add following to user
func (r *userRepo) AddFollowing(ctx context.Context, user *model.LocalUser, followingID string) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowingCollectionName)
	data := &model.RemoteUser{ID: followingID, CreatedAt: time.Now()}
	if err := addIfNotExists(ctx, col, data); err != nil {
		return err
	}
	return nil
}

// Remove follower from user
func (r *userRepo) RemoveFollower(ctx context.Context, user *model.LocalUser, followerID string) error {
	collection := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
	if err := removeIfExists(ctx, collection, followerID); err != nil {
		return err
	}
	return nil
}

// Remove following from user
func (r *userRepo) RemoveFollowing(ctx context.Context, user *model.LocalUser, followingID string) error {
	collection := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowingCollectionName)
	if err := removeIfExists(ctx, collection, followingID); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ListFollowers(ctx context.Context, user *model.LocalUser) ([]*model.RemoteUser, error) {
	users := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowersCollectionName)
	q := users.OrderBy("created_at", firestore.Desc)
	items, err := getAllItems[*model.RemoteUser](ctx, users, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *userRepo) ListFollowing(ctx context.Context, user *model.LocalUser) ([]*model.RemoteUser, error) {
	users := r.firestoreClient.Collection(UsersCollectionName).Doc(user.ID).Collection(FollowingCollectionName)
	q := users.OrderBy("created_at", firestore.Desc)
	items, err := getAllItems[*model.RemoteUser](ctx, users, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}
