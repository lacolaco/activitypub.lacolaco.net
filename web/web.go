package web

import (
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/gcp"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
	"github.com/lacolaco/activitypub.lacolaco.net/static"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/web/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/web/api"
	wellknown "github.com/lacolaco/activitypub.lacolaco.net/web/well-known"
	"go.uber.org/zap"
)

func Start(conf *config.Config) error {
	log.Print("starting server...")
	if conf.IsRunningOnCloud() {
		gin.SetMode(gin.ReleaseMode)
	}

	firestore := gcp.NewFirestoreClient()
	firebaseAuth := gcp.NewFirebaseAuthClient()
	defer firestore.Close()
	userRepo := repository.NewUserRepository(firestore)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(config.WithConfig(conf))
	r.Use(tracing.WithTracing(conf))
	r.Use(logging.WithLogging(conf))
	r.Use(auth.WithAuth(auth.FirebaseAuthTokenVerifier(firebaseAuth), userRepo.FindByUID))
	r.Use(static.WithStatic("/", "./public"))
	r.Use(errorHandler())
	r.Use(requestLogger())
	r.Use(func(ctx *gin.Context) {
		// set default cache-control header
		ctx.Header("Cache-Control", "no-cache")
		ctx.Next()
	})

	wellknown.New().Register(r)
	ap.New(userRepo).Register(r)
	api.New(userRepo).Register(r)

	// Start HTTP server.
	log.Printf("listening on http://localhost:%s", conf.Port)
	return r.Run(":" + conf.Port)
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logging.LoggerFromContext(c.Request.Context()).Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}

func requestLogger() gin.HandlerFunc {
	ignoredUserAgents := map[string]bool{
		"GoogleStackdriverMonitoring-UptimeChecks(https://cloud.google.com/monitoring)": true,
	}

	return func(c *gin.Context) {
		if ignoredUserAgents[c.Request.UserAgent()] {
			c.Next()
			return
		}
		logging.LoggerFromContext(c.Request.Context()).Debug("request.start",
			zap.String("method", c.Request.Method),
			zap.String("host", c.Request.Host),
			zap.String("url", c.Request.URL.String()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.String("accept", c.Request.Header.Get("Accept")),
			zap.Any("headers", c.Request.Header),
		)
		c.Next()
		logging.LoggerFromContext(c.Request.Context()).Debug("request.end",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.String("response.contentType", c.Writer.Header().Get("Content-Type")),
		)
	}
}
