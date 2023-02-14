package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type mockAuth struct {
	tokenToUID map[string]string
}

func (m *mockAuth) VerifyToken(ctx context.Context, idToken string) (auth.UID, error) {
	if uid, ok := m.tokenToUID[idToken]; ok {
		return uid, nil
	}
	return "", fmt.Errorf("invalid token")
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("do nothing if no authorization header", func(tt *testing.T) {
		mock := &mockAuth{}

		router := gin.New()
		router.Use(auth.WithAuth(mock.VerifyToken))
		router.GET("/test", func(c *gin.Context) {
			uid := auth.UIDFromContext(c.Request.Context())
			if uid != "" {
				tt.Error("uid is not empty")
			}
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("can get current user from valid authorization header", func(tt *testing.T) {
		user := &model.LocalUser{ID: "stub", UID: "test-uid"}
		mock := &mockAuth{
			tokenToUID: map[string]string{
				"token": user.UID,
			},
		}

		router := gin.New()
		router.Use(auth.WithAuth(mock.VerifyToken))
		router.GET("/test", func(c *gin.Context) {
			uid := auth.UIDFromContext(c.Request.Context())
			if uid == "" {
				tt.Error("uid is nil")
				return
			}
			if uid != user.UID {
				tt.Errorf("uid is not %s", user.UID)
			}
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestAssertAuthenticated(t *testing.T) {
	t.Run("throw 401 if no authorization header", func(tt *testing.T) {
		mock := &mockAuth{}
		router := gin.New()
		router.Use(auth.WithAuth(mock.VerifyToken))

		router.GET("/test", auth.AssertAuthenticated(), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("passthrough if request is authenticated", func(tt *testing.T) {
		user := &model.LocalUser{ID: "stub", UID: "test-uid"}
		mock := &mockAuth{
			tokenToUID: map[string]string{
				"token": user.UID,
			},
		}

		router := gin.New()
		router.Use(auth.WithAuth(mock.VerifyToken))
		router.GET("/test", auth.AssertAuthenticated(), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
