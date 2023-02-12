package middleware

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Static(prefix, dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		staticFS := static.LocalFile(dir, false)
		url := c.Request.URL.Path
		if staticFS.Exists(prefix, url) {
			c.FileFromFS(url, staticFS)
			return
		}
		c.Next()
		if c.Writer.Status() == http.StatusNotFound {
			c.Status(http.StatusOK)
			c.FileFromFS("/", staticFS)
		}
	}
}
