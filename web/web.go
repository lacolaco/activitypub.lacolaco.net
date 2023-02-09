package web

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	ap "github.com/lacolaco/activitypub.lacolaco.net/activitypub"
)

func Start(port string) error {
	r := gin.Default()

	r.GET("/.well-known/host-meta", handleWellKnownHostMeta)
	r.GET("/.well-known/webfinger", handleWebfinger)

	r.GET("/users/:username", handlePerson)
	r.GET("/@:username", handlePerson)
	r.GET("/", handler)

	// Start HTTP server.
	log.Printf("listening on http://localhost:%s", port)
	return r.Run(":" + port)
}

func handler(c *gin.Context) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	c.String(http.StatusOK, "Hello %s!", name)
}

func handlePerson(c *gin.Context) {
	username := c.Param("username")
	p := &ap.Person{
		Context:           "https://www.w3.org/ns/activitystreams",
		ID:                fmt.Sprintf("https://activitypub.lacolaco.net/users/%s", username),
		Name:              "lacolaco",
		PreferredUsername: username,
		Summary:           "I'm a software engineer.",
		Inbox:             fmt.Sprintf("https://activitypub.lacolaco.net/users/%s/inbox", username),
		Outbox:            fmt.Sprintf("https://activitypub.lacolaco.net/users/%s/outbox", username),
		URL:               fmt.Sprintf("https://activitypub.lacolaco.net/@%s", username),
		Icon: ap.Icon{
			Type:      "Image",
			MediaType: "image/png",
			URL:       "https://github.com/lacolaco.png",
		},
	}

	c.Header("Content-Type", "application/activity+json")
	c.JSON(http.StatusOK, p)
}

func handleWellKnownHostMeta(c *gin.Context) {
	c.Header("Content-Type", "application/xrd+xml")
	c.String(http.StatusOK, `<?xml version="1.0"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
	<Link rel="lrdd" type="application/xrd+xml" template="https://activitypub.lacolaco.net/.well-known/webfinger?resource={uri}" />
</XRD>
`)
}

func handleWebfinger(c *gin.Context) {
	resource := c.Query("resource")
	if resource == "" {
		c.String(http.StatusBadRequest, "resource is required")
		return
	}

	c.Header("Content-Type", "application/jrd+json")
	c.String(http.StatusOK, `{
	"subject": "acct:lacolaco@activitypub.lacolaco.net",
	"links": [
		{
			"rel": "self",
			"type": "application/activity+json",
			"href": "https://activitypub.lacolaco.net/users/lacolaco"
		}
	]
}`)
}
