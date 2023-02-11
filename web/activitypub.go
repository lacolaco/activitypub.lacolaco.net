package web

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	ap "github.com/lacolaco/activitypub.lacolaco.net/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
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

type activitypubEndpoints struct {
	userRepo UserRepository
}

func (e *activitypubEndpoints) RegisterRoutes(r *gin.Engine) {
	r.GET("/users/:username", e.handlePerson)
	r.GET("/@:username", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/users/"+ctx.Param("username"))
	})
	r.POST("/users/:username/inbox", e.handleInbox)
	r.GET("/users/:username/outbox", e.handleOutbox)
	r.GET("/users/:username/followers", e.handleFollowers)
	r.GET("/users/:username/following", e.handleFollowing)
	r.GET("/users/:username/collections/featured", e.handleFeatured)
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
	p := ap.NewPersonJSON(user, getBaseURI(c), conf.PublicKey)
	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, p)
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
	logger.Info("activity", zap.Any("type", activity.GetType()), zap.Any("actor", actor))

	switch activity.GetType() {
	case goap.FollowType:
		logger.Debug("accept follow", zap.String("from", string(actor.GetID())))
		err := s.userRepo.AddFollower(c.Request.Context(), user, string(actor.GetID()))
		if err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "add follower failed")
			return
		}
		if err := ap.Accept(c.Request.Context(), baseURI, user, activity); err != nil {
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
			if err := ap.Accept(c.Request.Context(), baseURI, user, activity); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}
			c.Status(http.StatusOK)
		}
	default:
		logger.Error("unsuppoted activity type")
		c.String(http.StatusBadRequest, "invalid activity type")
	}
}

func (s *activitypubEndpoints) handleOutbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res goap.Item
	isPage := c.Query("page") != ""
	if isPage {
		res = &goap.OrderedCollectionPage{
			ID:           goap.IRI(user.GetActivityPubID(baseURI) + "/outbox?page=true"),
			Type:         goap.OrderedCollectionPageType,
			TotalItems:   0,
			OrderedItems: []goap.Item{},
		}
	} else {
		res = &goap.OrderedCollection{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/outbox"),
			Type:       goap.OrderedCollectionType,
			TotalItems: 0,
			First:      goap.IRI(user.GetActivityPubID(baseURI) + "/outbox?page=true"),
		}
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func (s *activitypubEndpoints) handleFeatured(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res goap.Item
	isPage := c.Query("page") != ""
	if isPage {
		res = &goap.OrderedCollectionPage{
			ID:           goap.IRI(user.GetActivityPubID(baseURI) + "/collections/featured?page=true"),
			Type:         goap.OrderedCollectionPageType,
			TotalItems:   0,
			OrderedItems: []goap.Item{},
		}
	} else {
		res = &goap.OrderedCollection{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/collections/featured"),
			Type:       goap.OrderedCollectionType,
			TotalItems: 0,
			First:      goap.IRI(user.GetActivityPubID(baseURI) + "/collections/featured?page=true"),
		}
	}
	sendActivityJSON(c, http.StatusOK, res)
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
	users, err := s.userRepo.ListFollowers(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res goap.Item
	isPage := c.Query("page") != ""
	if isPage {
		res = &goap.OrderedCollectionPage{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/followers?page=true"),
			Type:       goap.OrderedCollectionPageType,
			TotalItems: uint(len(users)),
			OrderedItems: func() []goap.Item {
				items := make([]goap.Item, len(users))
				for i, item := range users {
					items[i] = goap.IRI(item.ID)
				}
				return items
			}(),
		}
	} else {
		res = &goap.OrderedCollection{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/followers"),
			Type:       goap.OrderedCollectionType,
			TotalItems: uint(len(users)),
			First:      goap.IRI(user.GetActivityPubID(baseURI) + "/followers?page=true"),
		}
	}
	sendActivityJSON(c, http.StatusOK, res)
}

func (s *activitypubEndpoints) handleFollowing(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	baseURI := getBaseURI(c)

	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	users, err := s.userRepo.ListFollowing(c.Request.Context(), user)
	if err != nil {
		logger.Error("failed to get followers", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res goap.Item
	isPage := c.Query("page") != ""
	if isPage {
		res = &goap.OrderedCollectionPage{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/following?page=true"),
			Type:       goap.OrderedCollectionPageType,
			TotalItems: uint(len(users)),
			OrderedItems: func() []goap.Item {
				items := make([]goap.Item, len(users))
				for i, item := range users {
					items[i] = goap.IRI(item.ID)
				}
				return items
			}(),
		}
	} else {
		res = &goap.OrderedCollection{
			ID:         goap.IRI(user.GetActivityPubID(baseURI) + "/following"),
			Type:       goap.OrderedCollectionType,
			TotalItems: uint(len(users)),
			First:      goap.IRI(user.GetActivityPubID(baseURI) + "/following?page=true"),
		}
	}
	sendActivityJSON(c, http.StatusOK, res)
}
