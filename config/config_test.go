package config_test

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"humungus.tedunangst.com/r/webs/httpsig"
)

func TestLoad(t *testing.T) {
	godotenv.Load("../.env")
	conf, err := config.Load()
	if err != nil {
		t.Errorf("got error, want no error: %v", err)
	}

	t.Run("PrivateKey is valid", func(tt *testing.T) {
		if conf.PrivateKey == nil {
			tt.Errorf("PrivateKey is nil")
		}
		_, err := httpsig.EncodeKey(conf.PrivateKey)
		if err != nil {
			tt.Errorf("PrivateKey key is not valid: %v", err)
		}
	})

	t.Run("PublicKey is valid", func(tt *testing.T) {
		if conf.PublicKey == nil {
			tt.Errorf("PublicKey is nil")
		}
		_, err := httpsig.EncodeKey(conf.PublicKey)
		if err != nil {
			tt.Errorf("PublicKey key is not valid: %v", err)
		}
	})
}
