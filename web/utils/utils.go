package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetBaseURI(c *gin.Context) string {
	return fmt.Sprintf("https://%s", c.Request.Host)
}
