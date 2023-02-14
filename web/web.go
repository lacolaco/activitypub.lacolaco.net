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
	"github.com/lacolaco/activitypub.lacolaco.net/web/hostmeta"
	"github.com/lacolaco/activitypub.lacolaco.net/web/webfinger"
)

func Start(conf *config.Config) error {
	log.Print("starting server...")
	if conf.IsRunningOnCloud() {
		gin.SetMode(gin.ReleaseMode)
	}

	closeTracer := tracing.InitTraceProvider(conf)
	defer closeTracer()

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
	r.Use(static.WithStatic("/", "./public"))
	r.Use(auth.WithAuth(auth.FirebaseAuthTokenVerifier(firebaseAuth)))
	r.Use(errorHandler())

	ap.New(userRepo).RegisterRoutes(r)
	api.New(userRepo).RegisterRoutes(r)
	hostmeta.RegisterRoutes(r)
	webfinger.RegisterRoutes(r)

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
