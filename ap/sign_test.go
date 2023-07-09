package ap_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
)

func TestSignRequest(t *testing.T) {
	godotenv.Load("../.env")
	conf, _ := config.Load()

	t.Run("can sign request", func(tt *testing.T) {
		body := []byte("hello")
		req, _ := http.NewRequest("POST", "https://example.com", bytes.NewReader(body))
		ap.SignRequest("keyID", conf.PrivateKey, req, body)
		defer func() {
			err := recover()
			if err != nil {
				tt.Fatal(err)
			}
		}()
	})
}
