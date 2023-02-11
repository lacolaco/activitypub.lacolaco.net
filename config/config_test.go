package config_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"humungus.tedunangst.com/r/webs/httpsig"
)

func TestLoad(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")
	conf, err := config.Load()
	if err != nil {
		t.Errorf("got error, want no error: %v", err)
	}

	t.Run("PrivateKey is valid", func(tt *testing.T) {
		if conf.PrivateKey == nil {
			tt.Errorf("PrivateKey is nil")
		}
		if conf.PrivateKey.Type != httpsig.RSA {
			tt.Errorf("PrivateKey type is not rsa")
		}
		_, err := httpsig.EncodeKey(conf.PrivateKey.Key)
		if err != nil {
			tt.Errorf("PrivateKey key is not valid: %v", err)
		}
	})

	t.Run("PublicKey is valid", func(tt *testing.T) {
		if conf.PublicKey == nil {
			tt.Errorf("PublicKey is nil")
		}
		if conf.PublicKey.Type != httpsig.RSA {
			tt.Errorf("PublicKey type is not rsa")
		}
		_, err := httpsig.EncodeKey(conf.PublicKey.Key)
		if err != nil {
			tt.Errorf("PublicKey key is not valid: %v", err)
		}
	})
}
