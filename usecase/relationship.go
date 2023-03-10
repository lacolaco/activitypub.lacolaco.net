package usecase

import (
	"context"
	"net/http"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/util"
)

type ErrMovedPermanently struct {
	NewURL string
}

func (e *ErrMovedPermanently) Error() string {
	return "moved permanently to " + e.NewURL
}

type UserRepository interface {
	FindByUID(ctx context.Context, uid model.UID) (*model.LocalUser, error)
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

func (u *relationshipUsecase) OnFollow(r *http.Request, uid model.UID, activity ap.ActivityObject) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnFollow")
	defer span.End()

	conf := config.ConfigFromContext(ctx)
	user, err := u.userRepo.FindByUID(ctx, uid)
	if err == repository.ErrNotFound {
		user, err = u.userRepo.FindByLocalID(ctx, string(uid))
		if err == nil {
			err = &ErrMovedPermanently{NewURL: util.GetBaseURI(r) + "/users/" + string(user.UID) + "/inbox"}
		}
	}
	if err != nil {
		return err
	}
	follower := model.NewFollower(string(activity.GetActor().GetID()), model.AttemptStatusCompleted)
	if err := u.userRepo.UpsertFollower(ctx, user, follower); err != nil {
		return err
	}
	actor := user.ToPerson(util.GetBaseURI(r), conf.PublicKey)
	if err := ap.Accept(ctx, actor, activity); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnUnfollow(r *http.Request, uid model.UID, activity ap.ActivityObject) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnUnfollow")
	defer span.End()

	conf := config.ConfigFromContext(ctx)
	user, err := u.userRepo.FindByUID(ctx, uid)
	if err == repository.ErrNotFound {
		user, err = u.userRepo.FindByLocalID(ctx, string(uid))
		if err == nil {
			err = &ErrMovedPermanently{NewURL: util.GetBaseURI(r) + "/users/" + string(user.UID) + "/inbox"}
		}
	}
	if err != nil {
		return err
	}
	if err := u.userRepo.DeleteFollower(ctx, user, string(activity.GetActor().GetID())); err != nil {
		return err
	}
	actor := user.ToPerson(util.GetBaseURI(r), conf.PublicKey)
	if err := ap.Accept(ctx, actor, activity); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnAcceptFollow(r *http.Request, uid model.UID, activity ap.ActivityObject) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnAcceptFollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
	if err == repository.ErrNotFound {
		user, err = u.userRepo.FindByLocalID(ctx, string(uid))
		if err == nil {
			err = &ErrMovedPermanently{NewURL: util.GetBaseURI(r) + "/users/" + string(user.UID) + "/inbox"}
		}
	}
	if err != nil {
		return err
	}
	following := model.NewFollowing(string(activity.GetActor().GetID()), model.AttemptStatusCompleted)
	if err := u.userRepo.UpsertFollowing(ctx, user, following); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) OnRejectFollow(r *http.Request, uid model.UID, activity ap.ActivityObject) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.OnRejectFollow")
	defer span.End()

	user, err := u.userRepo.FindByUID(ctx, uid)
	if err == repository.ErrNotFound {
		user, err = u.userRepo.FindByLocalID(ctx, string(uid))
		if err == nil {
			err = &ErrMovedPermanently{NewURL: util.GetBaseURI(r) + "/users/" + string(user.UID) + "/inbox"}
		}
	}
	if err != nil {
		return err
	}
	if err := u.userRepo.DeleteFollowing(ctx, user, string(activity.GetActor().GetID())); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) Follow(r *http.Request, uid model.UID, to string) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.Follow")
	defer span.End()

	conf := config.ConfigFromContext(ctx)
	user, err := u.userRepo.FindByUID(r.Context(), uid)
	if err != nil {
		return err
	}
	actor := user.ToPerson(util.GetBaseURI(r), conf.PublicKey)
	whom, err := ap.GetPerson(r.Context(), ap.IRI(to))
	if err != nil {
		return err
	}
	if err := ap.FollowPerson(r.Context(), actor, whom); err != nil {
		return err
	}
	following := model.NewFollowing(string(whom.GetID()), model.AttemptStatusPending)
	if err := u.userRepo.UpsertFollowing(r.Context(), user, following); err != nil {
		return err
	}
	return nil
}

func (u *relationshipUsecase) Unfollow(r *http.Request, uid model.UID, to string) error {
	ctx, span := tracing.StartSpan(r.Context(), "usecase.relationship.Unfollow")
	defer span.End()

	conf := config.ConfigFromContext(ctx)
	user, err := u.userRepo.FindByUID(r.Context(), uid)
	if err != nil {
		return err
	}
	actor := user.ToPerson(util.GetBaseURI(r), conf.PublicKey)
	whom, err := ap.GetPerson(r.Context(), ap.IRI(to))
	if err != nil {
		return err
	}
	if err := ap.UnfollowPerson(r.Context(), actor, whom); err != nil {
		return err
	}
	if err := u.userRepo.DeleteFollowing(r.Context(), user, string(whom.GetID())); err != nil {
		return err
	}
	return nil
}
