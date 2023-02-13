package webfinger_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/web/webfinger"
)

func TestHandler(t *testing.T) {

	t.Run("GET /.well-known/webfinger", func(tt *testing.T) {
		router := gin.New()

		webfinger.RegisterRoutes(router)
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
