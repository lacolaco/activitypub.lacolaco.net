package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/web/middleware"
)

func TestStatic(t *testing.T) {
	t.Run("can serve existing static file", func(tt *testing.T) {
		router := gin.New()
		router.Use(middleware.Static("/", "../fixtures/static"))
		req, _ := http.NewRequest("GET", "/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("can serve index.html with /", func(tt *testing.T) {
		router := gin.New()
		router.Use(middleware.Static("/", "../fixtures/static"))
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Result().Header.Get("Content-Type") != "text/html; charset=utf-8" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Content-Type"), "text/html; charset=utf-8")
		}
	})
}
