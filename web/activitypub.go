package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	ap "github.com/lacolaco/activitypub.lacolaco.net/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"go.uber.org/zap"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	AddFollower(ctx context.Context, user *model.User, followerID string) error
	RemoveFollower(ctx context.Context, user *model.User, followerID string) error
	ListFollowers(ctx context.Context, user *model.User) ([]*model.Follower, error)
}

type activitypubEndpoints struct {
	userRepo UserRepository
}

func (e *activitypubEndpoints) RegisterRoutes(r *gin.Engine) {
	r.GET("/users/:username", e.handlePerson)
	r.GET("/@:username", e.handlePerson)
	r.POST("/users/:username/inbox", e.handleInbox)
	r.GET("/users/:username/followers", e.handleFollowers)
}

func (s *activitypubEndpoints) handlePerson(c *gin.Context) {
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
	p := user.ToPerson(getBaseURI(c), &conf.RsaPrivateKey.PublicKey)
	sendActivityJSON(c, http.StatusOK, p)
}

func (s *activitypubEndpoints) handleInbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())

	if c.Request.Header.Get("Content-Type") != "application/activity+json" {
		logger.Sugar().Errorln("invalid content type", c.Request.Header.Get("Content-Type"))
		c.String(http.StatusBadRequest, "invalid content type")
		return
	}
	baseURI := getBaseURI(c)
	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	selfID := user.GetActivityPubID(baseURI)

	body, _ := io.ReadAll(c.Request.Body)
	o, err := goap.UnmarshalJSON(body)
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
	logger.Info("activity", zap.Any("activity", activity))
	from := activity.Actor

	switch activity.GetType() {
	case goap.FollowType:
		err := s.userRepo.AddFollower(c.Request.Context(), user, string(from.GetID()))
		if err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "add follower failed")
			return
		}
		logger.Debug("accept follow")
		actor, err := ap.GetActor(c.Request.Context(), string(from.GetID()))
		if err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "get actor failed")
			return
		}
		logger.Debug("post activity")
		res := ap.NewAccept(fmt.Sprintf("%s/%d", selfID, time.Now().Unix()), selfID, activity)
		if err := ap.PostActivity(c.Request.Context(), selfID, actor, res); err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "post activity failed")
			return
		}
		c.JSON(http.StatusOK, res)
	case goap.UndoType:
		switch activity.Object.GetType() {
		case goap.FollowType:
			err := s.userRepo.RemoveFollower(c.Request.Context(), user, string(from.GetID()))
			if err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}
			logger.Debug("accept follow", zap.String("from", string(from.GetID())))
			actor, err := ap.GetActor(c.Request.Context(), string(from.GetID()))
			if err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "invalid actor")
				return
			}
			logger.Debug("post activity to", zap.Any("actor", actor))
			res := ap.NewAccept(fmt.Sprintf("%s/%d", selfID, time.Now().Unix()), selfID, activity)
			if err := ap.PostActivity(c.Request.Context(), selfID, actor, res); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}

			c.JSON(http.StatusOK, res)
		}
	default:
		logger.Error("unsuppoted activity type")
		c.String(http.StatusBadRequest, "invalid activity type")
	}
}

func (s *activitypubEndpoints) handleFollowers(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	followers, err := s.userRepo.ListFollowers(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := &goap.OrderedCollection{
		Context:    goap.ActivityBaseURI,
		ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/followers"),
		Type:       goap.OrderedCollectionType,
		TotalItems: uint(len(followers)),
		OrderedItems: func() []goap.Item {
			items := make([]goap.Item, len(followers))
			for i, follower := range followers {
				items[i] = goap.IRI(follower.ID)
			}
			return items
		}(),
	}
	sendActivityJSON(c, http.StatusOK, res)
}
