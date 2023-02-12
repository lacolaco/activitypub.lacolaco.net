package ap

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/web/middleware"
	"go.uber.org/zap"
	"humungus.tedunangst.com/r/webs/httpsig"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	AddFollower(ctx context.Context, user *model.User, followerID string) error
	RemoveFollower(ctx context.Context, user *model.User, followerID string) error
	ListFollowers(ctx context.Context, user *model.User) ([]*model.RemoteUser, error)
	ListFollowing(ctx context.Context, user *model.User) ([]*model.RemoteUser, error)
}

type apService struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *apService {
	return &apService{
		userRepo: userRepo,
	}
}

func (s *apService) Register(r *gin.Engine) {
	userRoutes := r.Group("/users/:username", middleware.AssertAccept([]string{
		"application/activity+json",
		"application/ld+json",
		"application/json",
	}))

	userRoutes.GET("", s.handlePerson)
	userRoutes.POST("/inbox", middleware.AssertContentType([]string{"application/activity+json"}), s.handleInbox)
	userRoutes.GET("/outbox", s.handleOutbox)
	userRoutes.GET("/followers", s.handleFollowers)
	userRoutes.GET("/following", s.handleFollowing)
	r.GET("/@:username", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/users/"+ctx.Param("username"))
	})
}

func (s *apService) handlePerson(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	conf := config.FromContext(c.Request.Context())
	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	logger.Debug("user found", zap.Any("user", user))
	userPerson := ap.NewPerson(user, getBaseURI(c), conf.PublicKey)
	res := userPerson.AsMap()
	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, res)
}

func (s *apService) handleInbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	conf := config.FromContext(c.Request.Context())
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
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, getBaseURI(c), conf.PublicKey)

	switch activity.GetType() {
	case goap.FollowType:
		logger.Debug("accept follow", zap.String("from", string(actor.GetID())))
		err := s.userRepo.AddFollower(c.Request.Context(), user, string(actor.GetID()))
		if err != nil {
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
			logger.Debug("accept unfollow", zap.String("from", string(actor.GetID())))
			err := s.userRepo.RemoveFollower(c.Request.Context(), user, string(actor.GetID()))
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
		logger.Debug("create, update, delete", zap.Any("object", activity.Object))
		c.Status(200)
	case goap.AcceptType, goap.RejectType:
		logger.Debug("accept, reject", zap.Any("object", activity.Object))
		c.Status(200)
	default:
		logger.Error("unsuppoted activity type")
		c.String(http.StatusBadRequest, "invalid activity type")
	}
}

func (s *apService) handleOutbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	conf := config.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)
	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI, conf.PublicKey)
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
	conf := config.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI, conf.PublicKey)

	users, err := s.userRepo.ListFollowers(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &goap.OrderedCollection{
		ID:         goap.IRI(userPerson.FollowersURI()),
		Type:       goap.OrderedCollectionType,
		TotalItems: uint(len(users)),
		OrderedItems: func() []goap.Item {
			items := make([]goap.Item, len(users))
			for i, item := range users {
				items[i] = goap.IRI(item.ID)
			}
			return items
		}(),
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func (s *apService) handleFollowing(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	conf := config.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userPerson := ap.NewPerson(user, baseURI, conf.PublicKey)
	users, err := s.userRepo.ListFollowing(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &goap.OrderedCollection{
		ID:         goap.IRI(userPerson.FollowingURI()),
		Type:       goap.OrderedCollectionType,
		TotalItems: uint(len(users)),
		OrderedItems: func() []goap.Item {
			items := make([]goap.Item, len(users))
			for i, item := range users {
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

func getBaseURI(c *gin.Context) string {
	return fmt.Sprintf("https://%s", c.Request.Host)
}
