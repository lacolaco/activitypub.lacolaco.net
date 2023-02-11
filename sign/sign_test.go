package sign_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

func TestImportPrivateKey(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")

	t.Run("can import private key from string", func(tt *testing.T) {
		_, err := sign.DecodePrivateKey(os.Getenv("RSA_PRIVATE_KEY"))
		if err != nil {
			tt.Errorf("got error, want no error: %v", err)
		}
	})
}
