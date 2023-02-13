package static

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	ignoredPathPrefixes = []string{
		"/api",
		"/.well-known",
	}
)

func Serve(prefix, dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		staticFS := static.LocalFile(dir, false)
		url := c.Request.URL.Path
		for _, ignoredPrefix := range ignoredPathPrefixes {
			if strings.HasPrefix(url, ignoredPrefix) {
				return
			}
		}
		if staticFS.Exists(prefix, url) {
			if url == "/" {
				// index.html is not cacheable
				c.Header("Cache-Control", "no-cache")
			} else {
				c.Header("Cache-Control", "public, must-revalidate, max-age=0")
			}
			c.FileFromFS(url, staticFS)
			return
		}
		c.Next()
		if c.Writer.Status() == http.StatusNotFound {
			c.Header("Cache-Control", "no-cache")
			c.Status(http.StatusOK)
			c.FileFromFS("/", staticFS)
		}
	}
}
