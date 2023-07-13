package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r *userRepo) FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error) {
	collection := r.firestoreClient.Collection(UsersCollectionName)
	found, doc, err := findItem[model.LocalUser](ctx, collection.Where("id", "==", localID))
	if err != nil {
		return nil, err
	}
	found.UID = model.UID(doc.ID)
	return found, nil
}

func (r *userRepo) FindByUID(ctx context.Context, uid model.UID) (*model.LocalUser, error) {
	doc, err := r.firestoreClient.Collection(UsersCollectionName).Doc(string(uid)).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	user := &model.LocalUser{UID: model.UID(doc.Ref.ID)}
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// ===== Following =====

func (r *userRepo) UpsertFollowing(ctx context.Context, user *model.LocalUser, following *model.Following) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowingCollectionName)
	found, doc, err := findItem[model.Following](ctx, col.Where("user_id", "==", following.UserID))
	if err != nil && err != ErrNotFound {
		return err
	}
	if found != nil {
		following.CreatedAt = found.CreatedAt
	} else {
		doc = col.NewDoc()
		following.CreatedAt = time.Now()
	}
	if _, err := doc.Set(ctx, following); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ListFollowing(ctx context.Context, user *model.LocalUser) ([]*model.Following, error) {
	users := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowingCollectionName)
	q := users.OrderBy("created_at", firestore.Desc)
	items, err := getAllItems[model.Following](ctx, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *userRepo) FindFollowing(ctx context.Context, user *model.LocalUser, whom string) (*model.Following, error) {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowingCollectionName)
	found, doc, err := findItem[model.Following](ctx, col.Where("user_id", "==", whom))
	if err != nil {
		return nil, err
	}
	found.ID = doc.ID
	return found, nil
}

func (r *userRepo) DeleteFollowing(ctx context.Context, user *model.LocalUser, whom string) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowingCollectionName)
	if err := deleteItems(ctx, col.Where("user_id", "==", whom)); err != nil {
		return err
	}
	return nil
}

// ===== Followers =====

func (r *userRepo) UpsertFollower(ctx context.Context, user *model.LocalUser, follower *model.Follower) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowersCollectionName)
	found, doc, err := findItem[model.Follower](ctx, col.Where("id", "==", follower.ID))
	if err != nil && err != ErrNotFound {
		return err
	}
	if found != nil {
		follower.CreatedAt = found.CreatedAt
	} else {
		doc = col.NewDoc()
		follower.CreatedAt = time.Now()
	}
	if _, err := doc.Set(ctx, follower); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ListFollowers(ctx context.Context, user *model.LocalUser) ([]*model.Follower, error) {
	users := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowersCollectionName)
	q := users.OrderBy("created_at", firestore.Desc)
	items, err := getAllItems[model.Follower](ctx, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *userRepo) FindFollower(ctx context.Context, user *model.LocalUser, whom string) (*model.Follower, error) {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowersCollectionName)
	found, doc, err := findItem[model.Follower](ctx, col.Where("id", "==", whom))
	if err != nil {
		return nil, err
	}
	found.ID = doc.ID
	return found, nil
}

func (r *userRepo) DeleteFollower(ctx context.Context, user *model.LocalUser, whom string) error {
	col := r.firestoreClient.Collection(UsersCollectionName).Doc(user.GetDocID()).Collection(FollowersCollectionName)
	if err := deleteItems(ctx, col.Where("id", "==", whom)); err != nil {
		return err
	}
	return nil
}
