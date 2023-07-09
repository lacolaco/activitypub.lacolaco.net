package repository_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/gcp"
	"github.com/lacolaco/activitypub.lacolaco.net/repository"
)

func TestFindByUID(t *testing.T) {
	godotenv.Load("../.env")
	client := gcp.NewFirestoreClient()
	repo := repository.NewUserRepository(client)

	t.Run("can find user", func(tt *testing.T) {
		user, err := repo.FindByUID(context.Background(), "ug5pUwZAUGQBYhCPyHiCvWS5buc2")
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

	t.Run("can return ErrNotFound", func(tt *testing.T) {
		_, err := repo.FindByUID(context.Background(), "not-exist")
		if err != repository.ErrNotFound {
			tt.Fatal("err is not ErrNotFound")
		}
	})
}

func TestFindByLocalID(t *testing.T) {
	godotenv.Load("../.env")
	client := gcp.NewFirestoreClient()
	repo := repository.NewUserRepository(client)

	t.Run("can find user", func(tt *testing.T) {
		user, err := repo.FindByLocalID(context.Background(), "lacolaco")
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
