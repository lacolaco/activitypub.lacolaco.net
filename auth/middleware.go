package auth

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type contextKey int

const (
	currentUserContextKey contextKey = iota
)

type VerifyTokenFunc = func(ctx context.Context, token string) (UID, error)
type ResolveLocalUserFunc func(ctx context.Context, uid string) (*model.LocalUser, error)

func WithAuth(verifyToken VerifyTokenFunc, resolveLocalUser ResolveLocalUserFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.Header.Get("Authorization"), "Bearer ") {
			return
		}
		idToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		uid, err := verifyToken(c.Request.Context(), idToken)
		if err != nil {
			return
		}
		user, err := resolveLocalUser(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		ctx := context.WithValue(c.Request.Context(), currentUserContextKey, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ログイン中でなければ401を返す
func AssertAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := CurrentUserFromContext(c.Request.Context())
		if currentUser == nil {
			c.AbortWithStatus(401)
			return
		}
	}
}

// ログイン中のユーザーを取得する。ログイン中でなければnilを返す
func CurrentUserFromContext(c context.Context) *model.LocalUser {
	if user, ok := c.Value(currentUserContextKey).(*model.LocalUser); ok {
		return user
	}
	return nil
}
