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
	"github.com/lacolaco/activitypub.lacolaco.net/web/middleware"
	"github.com/lacolaco/activitypub.lacolaco.net/web/utils"
	"go.uber.org/zap"
	"humungus.tedunangst.com/r/webs/httpsig"
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

type apService struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *apService {
	return &apService{userRepo: userRepo}
}

func (s *apService) Register(r *gin.Engine) {
	assertJSONGet := middleware.AssertAccept([]string{
		"application/activity+json",
		"application/ld+json",
		"application/json",
	})
	assertJSONPost := middleware.AssertContentType([]string{"application/activity+json"})

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
	logger := logging.FromContext(c.Request.Context())
	conf := config.FromContext(c.Request.Context())
	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("user found", zap.Any("user", user))
	userPerson := ap.NewPerson(user, utils.GetBaseURI(c))
	res := userPerson.ToMap(conf.PublicKey)
	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, res)
}

func (s *apService) handleInbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	defer c.Request.Body.Close()
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("failed to get data", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("payload", zap.String("payload", string(payload)))
	if _, err := httpsig.VerifyRequest(c.Request, payload, httpsig.ActivityPubKeyGetter); err != nil {
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
	actor := activity.Actor
	logger.Debug("activity", zap.Any("type", activity.GetType()), zap.Any("actor", actor))
	username := c.Param("username")
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, utils.GetBaseURI(c))

	switch activity.GetType() {
	case goap.FollowType:
		follower := model.NewFollower(actor.GetID().String(), model.AttemptStatusCompleted)
		if err := s.userRepo.UpsertFollower(c.Request.Context(), user, follower); err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "add follower failed")
			return
		}
		if err := ap.Accept(c.Request.Context(), userPerson, activity); err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "post activity failed")
			return
		}
		c.Status(http.StatusOK)
	case goap.UndoType:
		switch activity.Object.GetType() {
		case goap.FollowType:
			err := s.userRepo.DeleteFollower(c.Request.Context(), user, actor.GetID().String())
			if err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}
			if err := ap.Accept(c.Request.Context(), userPerson, activity); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}
			c.Status(http.StatusOK)
		}
	case goap.CreateType, goap.UpdateType, goap.DeleteType:
		logger.Debug("not implemented: create, update, delete", zap.Any("object", activity.Object))
		c.Status(200)
	case goap.AnnounceType:
		logger.Debug("not implemented: announce", zap.Any("object", activity.Object))
		c.Status(200)
	case goap.AcceptType, goap.RejectType:
		switch {
		// follow request is accepted
		case activity.Object.GetType() == goap.FollowType && activity.GetType() == goap.AcceptType:
			following := model.NewFollowing(actor.GetID().String(), model.AttemptStatusCompleted)
			if err := s.userRepo.UpsertFollowing(c.Request.Context(), user, following); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "complete following failed")
				return
			}
			c.Status(200)
		// follow request is rejected
		case activity.Object.GetType() == goap.FollowType && activity.GetType() == goap.RejectType:
			if err := s.userRepo.DeleteFollowing(c.Request.Context(), user, actor.GetID().String()); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "cancel following failed")
				return
			}
			c.Status(200)
		// unfollow request is accepted
		case activity.Object.GetType() == goap.UndoType && activity.GetType() == goap.AcceptType:
			err := s.userRepo.DeleteFollowing(c.Request.Context(), user, actor.GetID().String())
			if err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "remove following failed")
				return
			}
			c.Status(200)
		// unfollow request is rejected
		case activity.Object.GetType() == goap.UndoType && activity.GetType() == goap.RejectType:
			logger.Info("unfollow request is rejected")
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
	logger := logging.FromContext(c.Request.Context())
	baseURI := utils.GetBaseURI(c)
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
	logger := logging.FromContext(c.Request.Context())
	baseURI := utils.GetBaseURI(c)

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
	logger := logging.FromContext(c.Request.Context())
	baseURI := utils.GetBaseURI(c)

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
