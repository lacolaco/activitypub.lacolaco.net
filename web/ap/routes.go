package ap

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/usecase"
	"github.com/lacolaco/activitypub.lacolaco.net/web/util"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type UserRepository interface {
	FindByLocalID(ctx context.Context, username string) (*model.LocalUser, error)
	UpsertFollower(ctx context.Context, user *model.LocalUser, follower *model.Follower) error
	ListFollowers(ctx context.Context, user *model.LocalUser) ([]*model.Follower, error)
	DeleteFollower(ctx context.Context, user *model.LocalUser, whom string) error
	UpsertFollowing(ctx context.Context, user *model.LocalUser, following *model.Following) error
	ListFollowing(ctx context.Context, user *model.LocalUser) ([]*model.Following, error)
	DeleteFollowing(ctx context.Context, user *model.LocalUser, whom string) error
}

type RelationshipUsecase interface {
	OnFollow(r *http.Request, username string, activity *goap.Activity) error
	OnUnfollow(r *http.Request, username string, activity *goap.Activity) error
	OnAcceptFollow(r *http.Request, username string, activity *goap.Activity) error
	OnRejectFollow(r *http.Request, username string, activity *goap.Activity) error
}

type apService struct {
	userRepo            UserRepository
	relationshipUsecase RelationshipUsecase
}

func New(userRepo UserRepository) *apService {
	return &apService{
		userRepo:            userRepo,
		relationshipUsecase: usecase.NewRelationshipUsecase(userRepo),
	}
}

func (s *apService) RegisterRoutes(r *gin.Engine) {
	assertJSONGet := util.AssertAccept([]string{
		"application/activity+json",
		"application/ld+json",
		"application/json",
	})
	assertJSONPost := util.AssertContentType([]string{"application/activity+json"})

	userRoutes := r.Group("/users/:username")
	userRoutes.GET("", assertJSONGet, s.handlePerson)
	userRoutes.POST("/inbox", assertJSONPost, s.handleInbox)
	userRoutes.GET("/outbox", assertJSONGet, s.handleOutbox)
	userRoutes.GET("/followers", assertJSONGet, s.handleFollowers)
	userRoutes.GET("/following", assertJSONGet, s.handleFollowing)
	r.GET("/@:username", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/users/"+ctx.Param("username"))
	})
}

func (s *apService) handlePerson(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	conf := config.ConfigFromContext(c.Request.Context())
	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("user found", zap.Any("user", user))
	userPerson := ap.NewPerson(user, util.GetBaseURI(c))
	res := userPerson.ToMap(conf.PublicKey)
	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, res)
}

func (s *apService) handleInbox(c *gin.Context) {
	ctx, span := tracing.StartSpan(c.Request.Context(), "ap.handleInbox")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	logger := logging.LoggerFromContext(ctx)
	defer c.Request.Body.Close()
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("failed to get data", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("payload", zap.String("payload", string(payload)))
	if err := ap.VerifyRequest(c.Request, payload); err != nil {
		logger.Error("failed to verify request", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	o, err := goap.UnmarshalJSON(payload)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var activity *goap.Activity
	err = goap.OnActivity(o, func(a *goap.Activity) error {
		activity = a
		return nil
	})
	if err != nil {
		logger.Error(err.Error())
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	span.SetAttributes(attribute.String("activity.actor", activity.Actor.GetID().String()))
	span.SetAttributes(attribute.String("activity.type", string(activity.GetType())))
	username := c.Param("username")

	switch activity.GetType() {
	case goap.FollowType:
		if err := s.relationshipUsecase.OnFollow(c.Request, username, activity); err != nil {
			c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	case goap.UndoType:
		switch activity.Object.GetType() {
		case goap.FollowType:
			if err := s.relationshipUsecase.OnUnfollow(c.Request, username, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(http.StatusOK)
		}
	case goap.CreateType, goap.UpdateType, goap.DeleteType:
		logger.Debug("not implemented: create, update, delete")
		c.Status(200)
	case goap.AnnounceType:
		logger.Debug("not implemented: announce", zap.Any("object", activity.Object))
		c.Status(200)
	case goap.AcceptType, goap.RejectType:
		switch {
		// follow request is accepted
		case activity.Object.GetType() == goap.FollowType && activity.GetType() == goap.AcceptType:
			if err := s.relationshipUsecase.OnAcceptFollow(c.Request, username, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(200)
		// follow request is rejected
		case activity.Object.GetType() == goap.FollowType && activity.GetType() == goap.RejectType:
			if err := s.relationshipUsecase.OnRejectFollow(c.Request, username, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(200)
		default:
			logger.Debug("not implemented: accept, reject", zap.Any("object", activity.Object))
			c.Status(200)
		}
	default:
		logger.Error("unsuppoted activity type")
		c.String(http.StatusBadRequest, "invalid activity type")
	}
}

func (s *apService) handleOutbox(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	baseURI := util.GetBaseURI(c)
	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI)
	res := &goap.OrderedCollection{
		ID:           goap.IRI(userPerson.OutboxURI()),
		Type:         goap.OrderedCollectionType,
		TotalItems:   0,
		OrderedItems: []goap.Item{},
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func (s *apService) handleFollowers(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	baseURI := util.GetBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI)
	followers, err := s.userRepo.ListFollowers(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &goap.OrderedCollection{
		ID:         goap.IRI(userPerson.FollowersURI()),
		Type:       goap.OrderedCollectionType,
		TotalItems: uint(len(followers)),
		OrderedItems: func() []goap.Item {
			items := make([]goap.Item, len(followers))
			for i, item := range followers {
				items[i] = goap.IRI(item.ID)
			}
			return items
		}(),
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func (s *apService) handleFollowing(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	baseURI := util.GetBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI)
	following, err := s.userRepo.ListFollowing(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &goap.OrderedCollection{
		ID:         goap.IRI(userPerson.FollowingURI()),
		Type:       goap.OrderedCollectionType,
		TotalItems: uint(len(following)),
		OrderedItems: func() []goap.Item {
			items := make([]goap.Item, len(following))
			for i, item := range following {
				items[i] = goap.IRI(item.ID)
			}
			return items
		}(),
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func sendActivityJSON(c *gin.Context, code int, item goap.Item) error {
	body, err := ap.MarshalActivityJSON(item)
	if err != nil {
		return err
	}
	c.Header("Content-Type", "application/activity+json")
	c.String(code, string(body))
	return nil
}
