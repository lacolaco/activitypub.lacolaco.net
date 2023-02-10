package web

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type wellKnownEndpoints struct{}

func (e *wellKnownEndpoints) RegisterRoutes(r *gin.Engine) {
	r.GET("/.well-known/webfinger", e.handleWebfinger)
}

func (e *wellKnownEndpoints) handleWebfinger(c *gin.Context) {
	host := c.Request.Host
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

	res := gin.H{
		"subject": "acct:" + username + "@" + host,
		"aliases": []string{
			"https://" + host + "/@" + username,
			"https://" + host + "/users/" + username,
		},
		"links": []interface{}{
			map[string]string{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + host + "/users/" + username,
			},
		},
	}
	c.Header("Content-Type", "application/jrd+json")
	c.JSON(http.StatusOK, res)
}
