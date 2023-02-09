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
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

type service struct {
	firestoreClient *firestore.Client
}

func Start(conf *config.Config) error {
	log.Print("starting server...")
	w := &service{
		firestoreClient: firestore.NewFirestoreClient(),
	}

	r := gin.Default()
	r.Use(config.Middleware(conf))
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
	username := c.Param("username")
	userDoc, err := s.firestoreClient.Collection("users").Doc(username).Get(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	user := &model.User{}
	if err := userDoc.DataTo(user); err != nil {
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
			PublicKeyPem: conf.RsaKeys.PublicKey,
		},
	}

	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, p)
}

func (s *service) handleInbox(c *gin.Context) {
	username := c.Param("username")
	if c.Request.Header.Get("Content-Type") != "application/activity+json" {
		fmt.Println("invalid content type", c.Request.Header.Get("Content-Type"))
		c.String(http.StatusBadRequest, "invalid content type")
		return
	}
	id := fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username)

	activity := &ap.Activity{}
	// log body json
	{
		req := c.Copy().Request
		body, _ := io.ReadAll(req.Body)
		fmt.Println("raw body")
		fmt.Printf("%s", string(body))
		o, err := goap.UnmarshalJSON(body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%#v", o)
	}
	if err := c.BindJSON(&activity); err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	fmt.Printf("%#v", activity)

	switch activity.Type {
	case ap.ActivityTypeFollow:
		followersCollection := s.firestoreClient.Collection("users").Doc(username).Collection("followers")
		_, err := followersCollection.Doc(activity.Actor.ID).Set(c.Request.Context(), map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		res := &ap.Activity{
			Context: ap.ActivityPubContext,
			Type:    ap.ActivityTypeAccept,
			Actor:   ap.Actor{ID: id},
			Object:  activity.ToObject(),
		}

		actor, err := ap.GetActor(c.Request.Context(), activity.Actor.ID)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "invalid actor")
			return
		}
		if err := ap.PostActivity(c.Request.Context(), id, actor, res); err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(http.StatusOK, res)
		return
	case ap.ActivityTypeUndo:
		switch activity.Object.Type {
		case ap.ActivityTypeFollow:
			// TODO: unfollow
			return
		}
	}

	fmt.Println("invalid activity type", activity.Type)
	c.String(http.StatusBadRequest, "invalid activity type")
}
