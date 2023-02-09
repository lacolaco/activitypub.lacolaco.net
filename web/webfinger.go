package web

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	sub, err := url.Parse(resource)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid resource")
		return
	}
	if sub.Scheme != "acct" {
		c.String(http.StatusBadRequest, "invalid resource")
		return
	}
	username := strings.Split(sub.Opaque, "@")[0]

	c.Header("Content-Type", "application/jrd+json")
	c.String(http.StatusOK, `{
	"subject": "acct:%s@activitypub.lacolaco.net",
	"links": [
		{
			"rel": "self",
			"type": "application/activity+json",
			"href": "https://activitypub.lacolaco.net/users/%s"
		}
	]
}`, username, username)
}
