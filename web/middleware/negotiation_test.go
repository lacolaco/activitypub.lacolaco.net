package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/web/middleware"
)

func TestAssertAccept(t *testing.T) {
	type args struct {
		expected []string
		input    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "match",
			args: args{expected: []string{"application/json"}, input: "application/json; charset=utf-8"},
			want: http.StatusOK,
		},
		{
			name: "mismatch",
			args: args{expected: []string{"application/json"}, input: "text/html; charset=utf-8"},
			want: http.StatusNotFound,
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(tt *testing.T) {
			router := gin.New()
			router.GET("/test", middleware.AssertAccept(spec.args.expected), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Accept", spec.args.input)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if spec.want != w.Code {
				tt.Errorf("got %d, want %d", w.Code, spec.want)
			}
		})
	}
}

func TestAssertContentType(t *testing.T) {
	type args struct {
		expected []string
		input    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "match",
			args: args{expected: []string{"application/json"}, input: "application/json; charset=utf-8"},
			want: http.StatusOK,
		},
		{
			name: "mismatch",
			args: args{expected: []string{"application/json"}, input: "text/plain; charset=utf-8"},
			want: http.StatusBadRequest,
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(tt *testing.T) {
			router := gin.New()
			router.POST("/test", middleware.AssertContentType(spec.args.expected), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			req, _ := http.NewRequest("POST", "/test", nil)
			req.Header.Set("Content-Type", spec.args.input)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if spec.want != w.Code {
				tt.Errorf("got %d, want %d", w.Code, spec.want)
			}
		})
	}
}
