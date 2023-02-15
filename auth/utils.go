package auth

import (
	"context"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

func FirebaseAuthTokenVerifier(client *firebaseauth.Client) VerifyTokenFunc {
	return func(ctx context.Context, token string) (model.UID, error) {
		t, err := client.VerifyIDToken(ctx, token)
		if err != nil {
			return "", err
		}
		return model.UID(t.UID), nil
	}
}
