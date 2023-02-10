package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	ap "github.com/lacolaco/activitypub.lacolaco.net/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	firestore "github.com/lacolaco/activitypub.lacolaco.net/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"go.uber.org/zap"
)

type service struct {
	firestoreClient *firestore.Client
}

func Start(conf *config.Config) error {
	log.Print("starting server...")
	w := &service{
		firestoreClient: firestore.NewFirestoreClient(),
	}
	if conf.IsRunningOnCloud() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(config.Middleware(conf))
	r.Use(config.Middleware(conf))
	r.Use(tracing.Middleware(conf))
	r.Use(logging.Middleware(conf))
	r.Use(errorHandler())
	r.Use(requestLogger())
	r.Use(func(ctx *gin.Context) {
		// set default cache-control header
		ctx.Header("Cache-Control", "no-store")
		ctx.Next()
	})

	r.GET("/.well-known/host-meta", handleWellKnownHostMeta)
	r.GET("/.well-known/webfinger", handleWebfinger)

	r.GET("/users/:username", w.handlePerson)
	r.GET("/@:username", w.handlePerson)
	r.POST("/users/:username/inbox", w.handleInbox)
	r.GET("/", w.handler)

	// Start HTTP server.
	log.Printf("listening on http://localhost:%s", conf.Port)
	return r.Run(":" + conf.Port)
}

func (s *service) handler(c *gin.Context) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	c.String(http.StatusOK, "Hello %s!", name)
}

func (s *service) handlePerson(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	username := c.Param("username")
	userDoc, err := s.firestoreClient.Collection("users").Doc(username).Get(c.Request.Context())
	if err != nil {
		logger.Error("failed to get user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	user := &model.User{}
	if err := userDoc.DataTo(user); err != nil {
		logger.Error("failed to parse user", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	conf := config.FromContext(c.Request.Context())

	id := fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username)
	p := &ap.Person{
		Context:           ap.ActivityPubContext,
		Type:              ap.ActivityTypePerson,
		ID:                id,
		Name:              user.Name,
		PreferredUsername: username,
		Summary:           user.Description,
		Inbox:             fmt.Sprintf("%s/inbox", id),
		Outbox:            fmt.Sprintf("%s/outbox", id),
		URL:               fmt.Sprintf("https://activitypub.lacolaco.net/@%s", username),
		Icon: ap.Icon{
			Type:      "Image",
			MediaType: user.Icon.MediaType,
			URL:       user.Icon.URL,
		},
		PublicKey: ap.PublicKey{
			Context:      ap.ActivityPubContext,
			Type:         "Key",
			ID:           fmt.Sprintf("%s#%s", id, sign.DefaultPublicKeyID),
			Owner:        id,
			PublicKeyPem: sign.ExportPublicKey(conf.RsaPrivateKey.PublicKey),
		},
	}

	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, p)
}

func (s *service) handleInbox(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	username := c.Param("username")
	if c.Request.Header.Get("Content-Type") != "application/activity+json" {
		logger.Sugar().Errorln("invalid content type", c.Request.Header.Get("Content-Type"))
		c.String(http.StatusBadRequest, "invalid content type")
		return
	}
	self := fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username)

	body, _ := io.ReadAll(c.Request.Body)
	logger.Sugar().Infoln("raw body")
	logger.Sugar().Infof("%s", string(body))
	o, err := goap.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	var activity *goap.Activity
	// log body json
	err = goap.OnActivity(o, func(a *goap.Activity) error {
		activity = a
		return nil
	})
	if err != nil {
		logger.Error(err.Error())
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	logger.Sugar().Infof("%#v", activity)
	from := activity.Actor

	switch activity.Type {
	case goap.FollowType:
		followersCollection := s.firestoreClient.Collection("users").Doc(username).Collection("followers")
		_, _, err := followersCollection.Add(c.Request.Context(), map[string]interface{}{
			"id": string(from.GetID()),
		})
		if err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		res := &goap.Accept{
			Context: activity.Context,
			Type:    goap.AcceptType,
			Actor:   goap.IRI(self),
			Object:  activity,
		}

		logger.Debug("accept follow", zap.String("from", string(from.GetID())))
		logger.Debug("get actor")
		actor, err := ap.GetActor(c.Request.Context(), string(from.GetID()))
		if err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "invalid actor")
			return
		}
		logger.Debug("actor", zap.String("actor", actor.ID))
		logger.Debug("post activity")
		if err := ap.PostActivity(c.Request.Context(), self, actor, res); err != nil {
			logger.Error(err.Error())
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(http.StatusOK, res)
		return
		// case ap.ActivityTypeUndo:
		// 	switch activity.Object.Type {
		// 	case ap.ActivityTypeFollow:
		// 		// TODO: unfollow
		// 		return
		// 	}
	}

	logger.Error("invalid activity type")
	c.String(http.StatusBadRequest, "invalid activity type")
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logging.FromContext(c.Request.Context()).Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.FromContext(c.Request.Context()).Debug("request.start",
			zap.String("method", c.Request.Method),
			zap.String("host", c.Request.Host),
			zap.String("url", c.Request.URL.String()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
		)
		c.Next()
		logging.FromContext(c.Request.Context()).Debug("request.end",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
