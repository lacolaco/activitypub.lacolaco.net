package auth

import (
	"context"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type contextKey int

const (
	currentUserContextKey contextKey = iota
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.LocalUser, error)
}

func Authenticate(authClient *auth.Client, userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.Header.Get("Authorization"), "Bearer ") {
			return
		}
		idToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			return
		}
		user, err := userRepo.FindByID(c.Request.Context(), token.UID)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		ctx := context.WithValue(c.Request.Context(), currentUserContextKey, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AssertAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := FromContext(c.Request.Context())
		if currentUser == nil {
			c.AbortWithStatus(401)
			return
		}
	}
}

func FromContext(c context.Context) *model.LocalUser {
	if user, ok := c.Value(currentUserContextKey).(*model.LocalUser); ok {
		return user
	}
	return nil
}
