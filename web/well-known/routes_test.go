package well_known_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	well_known "github.com/lacolaco/activitypub.lacolaco.net/web/well-known"
)

func TestWellKnown(t *testing.T) {
	router := gin.Default()
	ep := well_known.New()
	ep.Register(router)

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

	t.Run("GET /.well-known/webfinger", func(tt *testing.T) {
		req, err := http.NewRequest("GET", "https://lacolaco.example/.well-known/webfinger?resource=acct:alice%40lacolaco.example", nil)
		if err != nil {
			tt.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			tt.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var body map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			tt.Fatal(err)
		}
		if body["subject"] != "acct:alice@lacolaco.example" {
			tt.Errorf("unexpected subject: %v", body["subject"])
		}
		if body["aliases"] == nil {
			tt.Errorf("no aliases")
		}
		if body["links"] == nil {
			tt.Errorf("no links")
		}
	})
}
