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
	FindByUID(ctx context.Context, uid model.UID) (*model.LocalUser, error)
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

func (u *relationshipUsecase) OnFollow(r *http.Request, uid model.UID, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnFollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
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

func (u *relationshipUsecase) OnUnfollow(r *http.Request, uid model.UID, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnUnfollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
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

func (u *relationshipUsecase) OnAcceptFollow(r *http.Request, uid model.UID, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnAcceptFollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
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

func (u *relationshipUsecase) OnRejectFollow(r *http.Request, uid model.UID, activity *goap.Activity) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnRejectFollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
	if err != nil {
		return err
	}
	actor := activity.Actor
	if err := u.userRepo.DeleteFollowing(ctx, user, actor.GetID().String()); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) Follow(r *http.Request, uid model.UID, to string) error {
	user, err := u.userRepo.FindByUID(r.Context(), uid)
	if err != nil {
		return err
	}
	actor := ap.NewPerson(user, r.Host)
	whom, err := ap.GetPerson(r.Context(), to)
	if err != nil {
		return err
	}
	if err := ap.FollowPerson(r.Context(), actor, whom); err != nil {
		return err
	}
	following := model.NewFollowing(whom.GetID().String(), model.AttemptStatusPending)
	if err := u.userRepo.UpsertFollowing(r.Context(), user, following); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) Unfollow(r *http.Request, uid model.UID, to string) error {
	user, err := u.userRepo.FindByUID(r.Context(), uid)
	if err != nil {
		return err
	}
	actor := ap.NewPerson(user, r.Host)
	whom, err := ap.GetPerson(r.Context(), to)
	if err != nil {
		return err
	}
	if err := ap.UnfollowPerson(r.Context(), actor, whom); err != nil {
		return err
	}
	if err := u.userRepo.DeleteFollowing(r.Context(), user, whom.GetID().String()); err != nil {
		return err
	}
	return nil
}
