package hostmeta_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/web/hostmeta"
)

func TestWellKnown(t *testing.T) {
	router := gin.Default()
	hostmeta.RegisterRoutes(router)

	t.Run("GET /.well-known/host-meta", func(tt *testing.T) {
		req, err := http.NewRequest("GET", "https://lacolaco.example/.well-known/host-meta", nil)
		if err != nil {
			tt.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			tt.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
    <Link rel="lrdd" template="https://lacolaco.example/.well-known/webfinger?resource={uri}"/>
</XRD>`
		if strings.TrimSpace(rr.Body.String()) != expected {
			tt.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})

}
