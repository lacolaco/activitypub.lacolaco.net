package webfinger_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/web/webfinger"
)

type mockUserRepo struct {
	localIDToUser map[string]*model.LocalUser
}

func (r *mockUserRepo) FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error) {
	if user, ok := r.localIDToUser[localID]; ok {
		return user, nil
	}
	return nil, errors.New("not found")
}

func TestHandler(t *testing.T) {

	t.Run("GET /.well-known/webfinger", func(tt *testing.T) {
		router := gin.New()

		userRepo := &mockUserRepo{
			localIDToUser: map[string]*model.LocalUser{
				"alice": {
					ID:  "alice",
					UID: "alice-uid",
				},
			},
		}
		webfinger.New(userRepo).RegisterRoutes(router)
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
		if body["links"] == nil {
			tt.Errorf("no links")
		}
		// check links includes self and profile-page
		links := body["links"].([]interface{})
		if len(links) != 2 {
			tt.Errorf("unexpected links: %v", links)
		}
	})
}
