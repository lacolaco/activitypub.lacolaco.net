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
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"go.uber.org/zap"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	AddFollower(ctx context.Context, username, followerID string) error
	RemoveFollower(ctx context.Context, username, followerID string) error
}

type activitypubEndpoints struct {
	userRepo UserRepository
}

func (e *activitypubEndpoints) RegisterRoutes(r *gin.Engine) {
	r.GET("/users/:username", e.handlePerson)
	r.GET("/@:username", e.handlePerson)
	r.POST("/users/:username/inbox", e.handleInbox)
}

func (s *activitypubEndpoints) handlePerson(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	username := c.Param("username")
	user, err := s.userRepo.FindByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	conf := config.FromContext(c.Request.Context())

	id := fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username)
	p := &goap.Person{
		Context:           goap.ActivityBaseURI,
		Type:              goap.PersonType,
		ID:                goap.IRI(id),
		Name:              goap.DefaultNaturalLanguageValue(user.Name),
		PreferredUsername: goap.DefaultNaturalLanguageValue(username),
		Summary:           goap.DefaultNaturalLanguageValue(user.Description),
		Inbox:             goap.IRI(fmt.Sprintf("%s/inbox", id)),
		Outbox:            goap.IRI(fmt.Sprintf("%s/outbox", id)),
		URL:               goap.IRI(fmt.Sprintf("https://activitypub.lacolaco.net/@%s", username)),
		Icon: &goap.Object{
			Type:      "Image",
			MediaType: goap.MimeType(user.Icon.MediaType),
			URL:       goap.IRI(user.Icon.URL),
		},
		PublicKey: goap.PublicKey{
			ID:           goap.ID(fmt.Sprintf("%s#%s", id, sign.DefaultPublicKeyID)),
			Owner:        goap.IRI(id),
			PublicKeyPem: sign.ExportPublicKey(&conf.RsaPrivateKey.PublicKey),
		},
	}

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
	username := c.Param("username")
	self := fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username)

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
		err := s.userRepo.AddFollower(c.Request.Context(), username, string(from.GetID()))
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
		accept := ap.NewAccept(fmt.Sprintf("%s/%d", self, time.Now().Unix()), self, activity)
		if err := ap.PostActivity(c.Request.Context(), self, actor, accept); err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "post activity failed")
			return
		}
		c.JSON(http.StatusOK, accept)
	case goap.UndoType:
		switch activity.Object.GetType() {
		case goap.FollowType:
			err := s.userRepo.RemoveFollower(c.Request.Context(), username, string(from.GetID()))
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
			accept := ap.NewAccept(fmt.Sprintf("%s/%d", self, time.Now().Unix()), self, activity)
			if err := ap.PostActivity(c.Request.Context(), self, actor, accept); err != nil {
				logger.Error(err.Error())
				c.String(http.StatusInternalServerError, "internal server error")
				return
			}

			c.JSON(http.StatusOK, accept)
		}
	default:
		logger.Error("unsuppoted activity type")
		c.String(http.StatusBadRequest, "invalid activity type")
	}
}
