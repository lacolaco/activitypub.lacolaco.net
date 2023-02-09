package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

func GetActor(ctx context.Context, id string) (*Actor, error) {
	addr := fmt.Sprintf("https://%s", id)
	req, err := http.NewRequestWithContext(ctx, "GET", addr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/activity+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	// read actor from body
	actor := &Actor{}
	if err := json.NewDecoder(resp.Body).Decode(actor); err != nil {
		return nil, err
	}
	return actor, nil
}

func PostActivity(ctx context.Context, to *Actor, activity *Activity) error {
	addr := to.Inbox
	req, err := http.NewRequestWithContext(ctx, "POST", addr, nil)
	if err != nil {
		return err
	}
	signer, err := sign.NewHeaderSigner()
	if err != nil {
		return err
	}
	payload, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	conf := config.FromContext(ctx)
	req.Header.Set("Content-Type", "application/activity+json")
	signer.SignRequest(conf.RsaKeys.PrivateKey, sign.DefaultPublicKeyID, req, payload)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
