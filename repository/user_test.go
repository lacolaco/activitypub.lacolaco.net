package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/gcp"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
)

func TestFindByUsername(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")
	client := gcp.NewFirestoreClient()
	repo := repository.NewUserRepository(client)

	t.Run("can find user", func(tt *testing.T) {
		user, err := repo.FindByUsername(context.Background(), "lacolaco")
		if err != nil {
			tt.Fatal(err)
		}
		if user == nil {
			tt.Fatal("user is nil")
		}
		if user.ID != "lacolaco" {
			tt.Fatal("user.ID is not lacolaco")
		}
	})
}
