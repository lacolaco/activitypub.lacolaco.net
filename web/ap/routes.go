package ap

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/util"
	webutil "github.com/lacolaco/activitypub.lacolaco.net/web/util"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	mimeTypeActivityJSON = "application/activity+json"
)

type UserRepository interface {
	FindByUID(ctx context.Context, uid model.UID) (*model.LocalUser, error)
	ListFollowers(ctx context.Context, user *model.LocalUser) ([]*model.Follower, error)
	ListFollowing(ctx context.Context, user *model.LocalUser) ([]*model.Following, error)
}

type RelationshipUsecase interface {
	OnFollow(r *http.Request, uid model.UID, activity *ap.Activity) error
	OnUnfollow(r *http.Request, uid model.UID, activity *ap.Activity) error
	OnAcceptFollow(r *http.Request, uid model.UID, activity *ap.Activity) error
	OnRejectFollow(r *http.Request, uid model.UID, activity *ap.Activity) error
}

type apService struct {
	userRepo            UserRepository
	relationshipUsecase RelationshipUsecase
}

func New(userRepo UserRepository, relationshipUsecase RelationshipUsecase) *apService {
	return &apService{
		userRepo:            userRepo,
		relationshipUsecase: relationshipUsecase,
	}
}

func (s *apService) RegisterRoutes(r *gin.Engine) {
	assertJSONGet := webutil.AssertAccept([]string{
		"application/activity+json",
		"application/ld+json",
		"application/json",
	})
	assertJSONPost := webutil.AssertContentType([]string{"application/activity+json"})

	userRoutes := r.Group("/users/:uid")
	userRoutes.GET("", assertJSONGet, s.handlePerson)
	userRoutes.POST("/inbox", assertJSONPost, s.handleInbox)
	userRoutes.GET("/outbox", assertJSONGet, s.handleOutbox)
	userRoutes.GET("/followers", assertJSONGet, s.handleFollowers)
	userRoutes.GET("/following", assertJSONGet, s.handleFollowing)
}

func (s *apService) handlePerson(c *gin.Context) {
	conf := config.ConfigFromContext(c.Request.Context())
	uid := c.Param("uid")
	user, err := s.userRepo.FindByUID(c.Request.Context(), model.UID(uid))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	person := user.ToPerson(util.GetBaseURI(c.Request), conf.PublicKey)
	b, err := person.ToBytes()
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, mimeTypeActivityJSON, b)
}

func (s *apService) handleInbox(c *gin.Context) {
	ctx, span := tracing.StartSpan(c.Request.Context(), "ap.handleInbox")
	defer span.End()
	c.Request = c.Request.WithContext(ctx)

	logger := logging.LoggerFromContext(ctx)
	defer c.Request.Body.Close()
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("failed to get data", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("payload", zap.String("payload", string(b)))
	if err := ap.VerifyRequest(c.Request, b); err != nil {
		logger.Error("failed to verify request", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	activity, err := ap.ActivityFromBytes(b)
	if err != nil {
		c.Error(err)
		return
	}
	span.SetAttributes(attribute.String("activity.actor", activity.Actor))
	span.SetAttributes(attribute.String("activity.type", string(activity.Type)))
	uid := model.UID(c.Param("uid"))

	switch activity.Type {
	case ap.FollowActivityType:
		if err := s.relationshipUsecase.OnFollow(c.Request, uid, activity); err != nil {
			c.Error(err)
			return
		}
		c.Status(http.StatusOK)
	case ap.UndoActivityType:
		switch activity.Object.GetType() {
		case string(ap.FollowActivityType):
			if err := s.relationshipUsecase.OnUnfollow(c.Request, uid, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(http.StatusOK)
		}
	case ap.CreateActivityType, ap.UpdateActivityType, ap.DeleteActivityType:
		logger.Debug("not implemented: create, update, delete")
		c.Status(200)
	case ap.AnnounceActivityType:
		logger.Debug("not implemented: announce", zap.Any("object", activity.Object))
		c.Status(200)
	case ap.AcceptActivityType, ap.RejectActivityType:
		switch {
		// follow request is accepted
		case activity.GetType() == string(ap.AcceptActivityType) && activity.Object.GetType() == string(ap.FollowActivityType):
			if err := s.relationshipUsecase.OnAcceptFollow(c.Request, uid, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(200)
		// follow request is rejected
		case activity.GetType() == string(ap.RejectActivityType) && activity.Object.GetType() == string(ap.FollowActivityType):
			if err := s.relationshipUsecase.OnRejectFollow(c.Request, uid, activity); err != nil {
				c.Error(err)
				return
			}
			c.Status(200)
		default:
			logger.Debug("not implemented: accept, reject")
			c.Status(200)
		}
	default:
		logger.Debug("unsuppoted activity type")
		c.Status(http.StatusOK)
	}
}

func (s *apService) handleOutbox(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	conf := config.ConfigFromContext(c.Request.Context())
	uid := c.Param("uid")
	user, err := s.userRepo.FindByUID(c.Request.Context(), model.UID(uid))
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	person := user.ToPerson(util.GetBaseURI(c.Request), conf.PublicKey)
	res := &ap.OrderedCollection{
		ID:           person.Outbox,
		TotalItems:   0,
		OrderedItems: []ap.ActivityPubObject{},
	}
	b, err := res.ToBytes()
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, mimeTypeActivityJSON, b)
}

func (s *apService) handleFollowers(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	conf := config.ConfigFromContext(c.Request.Context())

	uid := c.Param("uid")
	user, err := s.userRepo.FindByUID(c.Request.Context(), model.UID(uid))
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	person := user.ToPerson(util.GetBaseURI(c.Request), conf.PublicKey)
	followers, err := s.userRepo.ListFollowers(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &ap.OrderedCollection{
		ID:         person.Followers,
		TotalItems: len(followers),
		OrderedItems: func() []ap.ActivityPubObject {
			items := make([]ap.ActivityPubObject, len(followers))
			for i, item := range followers {
				items[i] = ap.IRI(item.UserID)
			}
			return items
		}(),
	}
	b, err := res.ToBytes()
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, mimeTypeActivityJSON, b)
}

func (s *apService) handleFollowing(c *gin.Context) {
	logger := logging.LoggerFromContext(c.Request.Context())
	conf := config.ConfigFromContext(c.Request.Context())

	uid := c.Param("uid")
	user, err := s.userRepo.FindByUID(c.Request.Context(), model.UID(uid))
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	person := user.ToPerson(util.GetBaseURI(c.Request), conf.PublicKey)
	following, err := s.userRepo.ListFollowing(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &ap.OrderedCollection{
		ID:         person.Following,
		TotalItems: len(following),
		OrderedItems: func() []ap.ActivityPubObject {
			items := make([]ap.ActivityPubObject, len(following))
			for i, item := range following {
				items[i] = ap.IRI(item.UserID)
			}
			return items
		}(),
	}
	b, err := res.ToBytes()
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, mimeTypeActivityJSON, b)
}
