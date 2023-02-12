package ap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"go.uber.org/zap"
	"humungus.tedunangst.com/r/webs/httpsig"
)

var (
	HTTPClient = http.DefaultClient
)

const (
	userAgent            = "activitypub.lacolaco.net/1.0"
	mimeTypeActivityJSON = "application/activity+json"
)

func getActor(ctx context.Context, publicKeyID string, id string) (*goap.Actor, error) {
	addr, _ := url.Parse(id)
	if addr.Scheme == "" {
		addr.Scheme = "https"
	}
	body, err := getActivityJSON(ctx, "", addr.String())
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

func getActivityJSON(ctx context.Context, publicKeyID string, url string) ([]byte, error) {
	conf := config.FromContext(ctx)
	logger := logging.FromContext(ctx)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mimeTypeActivityJSON)
	req.Header.Set("User-Agent", userAgent)
	key := httpsig.PrivateKey{Key: conf.PrivateKey, Type: httpsig.RSA}
	httpsig.SignRequest(publicKeyID, key, req, nil)
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req = req.WithContext(c)
	logger.Debug("getActivityJSON.request", zap.String("url", req.URL.String()), zap.Any("headers", req.Header))
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Debug("getActivityJSON.response", zap.Any("response", string(body)))
	return body, nil
}

func postActivityJSON(ctx context.Context, publicKeyID string, url string, body []byte) ([]byte, error) {
	conf := config.FromContext(ctx)
	logger := logging.FromContext(ctx)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", mimeTypeActivityJSON)
	key := httpsig.PrivateKey{Key: conf.PrivateKey, Type: httpsig.RSA}
	httpsig.SignRequest(publicKeyID, key, req, body)
	logger.Debug("postActivityJSON.request", zap.Any("headers", req.Header))
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req = req.WithContext(c)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("postActivityJSON.response", zap.Any("response", string(respBody)))
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusCreated:
	default:
		return nil, fmt.Errorf("http post status: %d", resp.StatusCode)
	}
	return respBody, nil
}
