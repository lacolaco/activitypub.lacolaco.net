package static_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/static"
)

func TestStaticMiddleware(t *testing.T) {
	t.Run("can serve existing static file", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Result().Header.Get("Cache-Control") != "no-cache" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Cache-Control"), "no-cache")
		}
	})

	t.Run("can serve index.html with /", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Result().Header.Get("Content-Type") != "text/html; charset=utf-8" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Content-Type"), "text/html; charset=utf-8")
		}
		if w.Result().Header.Get("Cache-Control") != "no-cache" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Cache-Control"), "no-cache")
		}
	})

	t.Run("can serve unknown path with /", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/unknown", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Result().Header.Get("Content-Type") != "text/html; charset=utf-8" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Content-Type"), "text/html; charset=utf-8")
		}
		if w.Result().Header.Get("Cache-Control") != "no-cache" {
			tt.Errorf("got %s, want %s", w.Result().Header.Get("Cache-Control"), "no-cache")
		}
	})

	t.Run("respect pre-defined handlers", func(tt *testing.T) {
		router := gin.New()
		router.GET("/test.txt", func(c *gin.Context) {
			c.String(http.StatusOK, "from handler")
		})
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.String() != "from handler" {
			tt.Errorf("got %s, want %s", w.Body.String(), "from handler")
		}
	})

	t.Run("skip post-defined handler", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		router.GET("/test.txt", func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		})
		req, _ := http.NewRequest("GET", "/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("skip post request", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("POST", "/test.txt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("ignore api routes", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/api/foo", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("ignore .well-known routes", func(tt *testing.T) {
		router := gin.New()
		router.Use(static.WithStatic("/", "./fixtures/static"))
		req, _ := http.NewRequest("GET", "/.well-known/foo", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			tt.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}
