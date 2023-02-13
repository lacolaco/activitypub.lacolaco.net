package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AssertAccept(expected []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		for _, e := range expected {
			if strings.Contains(accept, e) {
				return
			}
		}
		c.Status(http.StatusNotFound)
		c.Abort()
	}
}

func AssertContentType(expected []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.Request.Header.Get("Content-Type")
		for _, e := range expected {
			if strings.Contains(contentType, e) {
				return
			}
		}
		c.Status(http.StatusBadRequest)
		c.Abort()
	}
}
