package usecase

import (
	"context"
	"net/http"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
)

type UserRepository interface {
	FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error)
	UpsertFollowing(ctx context.Context, user *model.LocalUser, following *model.Following) error
	DeleteFollowing(ctx context.Context, user *model.LocalUser, whom string) error
	UpsertFollower(ctx context.Context, user *model.LocalUser, follower *model.Follower) error
	DeleteFollower(ctx context.Context, user *model.LocalUser, whom string) error
}

type relationshipUsecase struct {
	userRepo UserRepository
}

func NewRelationshipUsecase(userRepo UserRepository) *relationshipUsecase {
	return &relationshipUsecase{userRepo: userRepo}
}

func (u *relationshipUsecase) OnFollow(r *http.Request, username string, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnFollow")
	defer span.End()

	user, err := u.userRepo.FindByLocalID(ctx, username)
	if err != nil {
		return err
	}
	actor := activity.Actor
	follower := model.NewFollower(actor.GetID().String(), model.AttemptStatusCompleted)
	if err := u.userRepo.UpsertFollower(ctx, user, follower); err != nil {
		return err
	}
	acceptActor := ap.NewPerson(user, r.Host)
	if err := ap.Accept(ctx, acceptActor, activity); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnUnfollow(r *http.Request, username string, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnUnfollow")
	defer span.End()

	user, err := u.userRepo.FindByLocalID(ctx, username)
	if err != nil {
		return err
	}
	actor := activity.Actor
	if err := u.userRepo.DeleteFollower(ctx, user, actor.GetID().String()); err != nil {
		return err
	}
	acceptActor := ap.NewPerson(user, r.Host)
	if err := ap.Accept(ctx, acceptActor, activity); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnAcceptFollow(r *http.Request, username string, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnAcceptFollow")
	defer span.End()

	user, err := u.userRepo.FindByLocalID(ctx, username)
	if err != nil {
		return err
	}
	actor := activity.Actor
	following := model.NewFollowing(actor.GetID().String(), model.AttemptStatusCompleted)
	if err := u.userRepo.UpsertFollowing(ctx, user, following); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnRejectFollow(r *http.Request, username string, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnRejectFollow")
	defer span.End()

	user, err := u.userRepo.FindByLocalID(ctx, username)
	if err != nil {
		return err
	}
	actor := activity.Actor
	if err := u.userRepo.DeleteFollowing(ctx, user, actor.GetID().String()); err != nil {
		return err
	}
	return nil
}
