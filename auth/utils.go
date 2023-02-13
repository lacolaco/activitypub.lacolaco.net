package auth

import (
	"context"

	firebaseauth "firebase.google.com/go/v4/auth"
)

func FirebaseAuthTokenVerifier(client *firebaseauth.Client) VerifyTokenFunc {
	return func(ctx context.Context, token string) (UID, error) {
		t, err := client.VerifyIDToken(ctx, token)
		if err != nil {
			return "", err
		}
		return UID(t.UID), nil
	}
}
