package util

import (
	"fmt"
	"net/http"
)

func GetBaseURI(r *http.Request) string {
	return fmt.Sprintf("https://%s", r.Host)
}
