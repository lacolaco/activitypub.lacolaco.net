package sign_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

func TestImportPrivateKey(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")

	t.Run("can import private key from string", func(tt *testing.T) {
		_, err := sign.ImportPrivateKey(os.Getenv("RSA_PRIVATE_KEY"))
		if err != nil {
			tt.Errorf("got error, want no error: %v", err)
		}
	})
}

func TestExportPublicKey(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")
	conf, err := config.Load()
	if err != nil {
		t.Fatalf("got error, want no error: %v", err)
	}
	t.Logf("%#v", conf)

	t.Run("can export public key as string", func(tt *testing.T) {
		got := sign.ExportPublicKey(&conf.RsaPrivateKey.PublicKey)
		if got == "" {
			t.Errorf("got empty string, want public key")
		}
	})
}
