package static

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ignoredPathPrefixes = []string{
		"/api",
		"/.well-known",
	}
	cacheControl = "no-cache"
)

// WithStatic は静的ファイルを配信するミドルウェアを生成する
//
// パスと一致するファイルが存在する場合は、そのファイルを返す。
// パスと一致するファイルが存在しない場合は、次のハンドラを呼び出す。
// 後続のハンドラで404が返された場合は、index.htmlを返す
func WithStatic(prefix, dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodGet && method != http.MethodHead {
			c.Next()
			return
		}

		reqPath := c.Request.URL.Path
		for _, ignoredPrefix := range ignoredPathPrefixes {
			if strings.HasPrefix(reqPath, ignoredPrefix) {
				return
			}
		}
		if reqPath == "/" {
			reqPath = "/index.html"
		}
		filepath := path.Join(dir, strings.TrimPrefix(reqPath, prefix))
		if exists(filepath) {
			c.Header("Cache-Control", cacheControl)
			c.File(filepath)
			c.Abort()
			return
		}

		c.Next()

		if c.Writer.Status() == http.StatusNotFound {
			c.Header("Cache-Control", cacheControl)
			c.Status(http.StatusOK)
			c.File(path.Join(dir, "index.html"))
		}
	}
}

func exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
