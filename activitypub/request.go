package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

func GetActor(ctx context.Context, id string) (*Actor, error) {
	addr, _ := url.Parse(id)
	if addr.Scheme == "" {
		addr.Scheme = "https"
	}
	req, err := http.NewRequestWithContext(ctx, "GET", addr.String(), nil)
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

func PostActivity(ctx context.Context, from string, to *Actor, activity *goap.Activity) error {
	logger := logging.FromContext(ctx)
	addr := to.Inbox
	req, err := http.NewRequestWithContext(ctx, "POST", addr, nil)
	if err != nil {
		return err
	}
	signer, err := sign.NewHeaderSigner()
	if err != nil {
		return err
	}
	payload, err := activity.MarshalJSON()
	if err != nil {
		return err
	}

	conf := config.FromContext(ctx)
	req.Header.Set("Content-Type", "application/activity+json")
	keyId := fmt.Sprintf("%s#%s", from, sign.DefaultPublicKeyID)
	signer.SignRequest(conf.RsaPrivateKey, keyId, req, payload)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logger.Sugar().Infoln("raw body")
	logger.Sugar().Infof("%s", string(body))
	return nil
}
