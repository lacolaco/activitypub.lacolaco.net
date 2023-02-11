package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	goap "github.com/go-ap/activitypub"
	ap "github.com/lacolaco/activitypub.lacolaco.net/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	firestore "github.com/lacolaco/activitypub.lacolaco.net/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"go.uber.org/zap"
)

func Start(conf *config.Config) error {
	log.Print("starting server...")
	if conf.IsRunningOnCloud() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(config.Middleware(conf))
	r.Use(config.Middleware(conf))
	r.Use(tracing.Middleware(conf))
	r.Use(logging.Middleware(conf))
	r.Use(errorHandler())
	r.Use(requestLogger())
	r.Use(func(ctx *gin.Context) {
		// set default cache-control header
		ctx.Header("Cache-Control", "public, no-cache")
		ctx.Next()
	})

	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/index.html", "./static/index.html")
	r.StaticFile("/robots.txt", "./static/robots.txt")

	wkService := &wellKnownEndpoints{}
	wkService.RegisterRoutes(r)

	firestoreClient := firestore.NewFirestoreClient()
	apService := &activitypubEndpoints{
		userRepo: repository.NewUserRepository(firestoreClient),
	}
	apService.RegisterRoutes(r)

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
		logging.FromContext(c.Request.Context()).Error(err.Error())
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
		logging.FromContext(c.Request.Context()).Debug("request.start",
			zap.String("method", c.Request.Method),
			zap.String("host", c.Request.Host),
			zap.String("url", c.Request.URL.String()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.String("accept", c.Request.Header.Get("Accept")),
			zap.Any("headers", c.Request.Header),
		)
		if c.Request.Method == "POST" {
			var body map[string]interface{}
			c.ShouldBindJSON(&body)
			logging.FromContext(c.Request.Context()).Debug("request.body", zap.Any("body", body))
		}
		c.Next()
		logging.FromContext(c.Request.Context()).Debug("request.end",
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
			zap.String("response.contentType", c.Writer.Header().Get("Content-Type")),
		)
	}
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
