package webfinger

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type UserRepository interface {
	FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error)
}

type service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) RegisterRoutes(r *gin.Engine) {
	r.GET("/.well-known/webfinger", s.handle)
}

func (s *service) handle(c *gin.Context) {
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
	localID := strings.Split(sub.Opaque, "@")[0]
	user, err := s.userRepo.FindByLocalID(c, localID)
	if err != nil {
		c.String(http.StatusNotFound, "user not found")
		return
	}

	res := gin.H{
		"subject": "acct:" + localID + "@" + host,
		"links": []interface{}{
			map[string]string{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + host + "/users/" + string(user.UID),
			},
		},
	}
	c.Header("Content-Type", "application/jrd+json")
	c.Header("Cache-Control", "max-age=3600, public")
	c.JSON(http.StatusOK, res)
}
