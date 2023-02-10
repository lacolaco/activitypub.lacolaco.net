package activitypub

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"go.uber.org/zap"
)

func GetActor(ctx context.Context, id string) (*goap.Actor, error) {
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	item, err := goap.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	var actor *goap.Actor
	err = goap.OnActor(item, func(a *goap.Actor) error {
		actor = a
		return nil
	})
	if err != nil {
		return nil, err
	}
	return actor, nil
}

func PostActivity(ctx context.Context, from string, to *goap.Actor, activity *goap.Activity) error {
	conf := config.FromContext(ctx)
	logger := logging.FromContext(ctx)

	addr := string(to.Inbox.GetLink())
	payload, err := activity.MarshalJSON()
	if err != nil {
		return err
	}
	logger.Debug("raw payload", zap.Any("payload", string(payload)))

	req, err := http.NewRequestWithContext(ctx, "POST", addr, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	keyId := fmt.Sprintf("%s#%s", from, sign.DefaultPublicKeyID)
	signedHeaders, err := sign.SignedHeaders(payload, addr, conf.RsaPrivateKey, keyId)
	if err != nil {
		return err
	}
	for k, v := range signedHeaders {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logger.Debug("raw response", zap.Any("response", string(body)))
	return nil
}
