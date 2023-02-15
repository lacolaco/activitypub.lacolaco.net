package auth

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"go.opentelemetry.io/otel/attribute"
)

type contextKey int

const (
	uidContextKey contextKey = iota
)

type VerifyTokenFunc = func(ctx context.Context, token string) (model.UID, error)

func WithAuth(verifyToken VerifyTokenFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartSpan(c.Request.Context(), "auth.WithAuth")
		defer span.End()
		if !strings.HasPrefix(c.Request.Header.Get("Authorization"), "Bearer ") {
			span.SetAttributes(attribute.String("auth.uid", ""))
			c.Next()
			return
		}
		idToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
		uid, err := verifyToken(ctx, idToken)
		if err != nil {
			span.SetAttributes(attribute.String("auth.uid", ""))
			c.Next()
			return
		}
		span.SetAttributes(attribute.String("auth.uid", string(uid)))
		ctx = context.WithValue(ctx, uidContextKey, uid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ログイン中でなければ401を返す
func AssertAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := UIDFromContext(c.Request.Context())
		if uid == "" {
			c.AbortWithStatus(401)
			return
		}
	}
}

// ログイン中のユーザーを取得する。ログイン中でなければ空文字列を返す
func UIDFromContext(c context.Context) model.UID {
	if uid, ok := c.Value(uidContextKey).(model.UID); ok {
		return uid
	}
	return ""
}
