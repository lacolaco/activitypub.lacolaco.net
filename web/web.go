package web

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/gcp"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/usecase"
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
	relationshipUsecase := usecase.NewRelationshipUsecase(userRepo)
	searchUsecase := usecase.NewSearchUsecase()

	r := gin.New()
	r.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{conf.ClientOrigin}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(config.WithConfig(conf))
	r.Use(tracing.WithTracing(conf))
	r.Use(logging.WithLogging(conf))
	r.Use(auth.WithAuth(auth.FirebaseAuthTokenVerifier(firebaseAuth)))
	r.Use(errorHandler())

	ap.New(userRepo, relationshipUsecase).RegisterRoutes(r)
	api.New(userRepo, relationshipUsecase, searchUsecase).RegisterRoutes(r)
	hostmeta.RegisterRoutes(r)
	webfinger.New(userRepo).RegisterRoutes(r)

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
		if err.Err == repository.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, err.JSON())
			return
		}
		if moved, ok := err.Err.(*usecase.ErrMovedPermanently); ok {
			c.Redirect(http.StatusMovedPermanently, moved.NewURL)
			return
		}
		logging.LoggerFromContext(c.Request.Context()).Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}
