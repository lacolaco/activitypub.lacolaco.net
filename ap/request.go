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
)

var (
	HTTPClient = http.DefaultClient
)

const (
	userAgent            = "activitypub.lacolaco.net/1.0"
	mimeTypeActivityJSON = "application/activity+json"
	systemActorID        = "https://activitypub.lacolaco.net/users/system"
)

func getPerson(ctx context.Context, id string) (*goap.Actor, error) {
	addr, _ := url.Parse(id)
	if addr.Scheme == "" {
		addr.Scheme = "https"
	}
	body, err := getActivityJSON(ctx, systemActor, addr.String())
	if err != nil {
		return nil, err
	}
	item, err := goap.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	var p *goap.Person
	err = goap.OnActor(item, func(a *goap.Actor) error {
		p = a
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getActivityJSON(ctx context.Context, actor Actor, url string) ([]byte, error) {
	conf := config.FromContext(ctx)
	logger := logging.FromContext(ctx)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mimeTypeActivityJSON)
	req.Header.Set("User-Agent", userAgent)
	publicKeyID := GetPublicKeyID(actor)
	SignRequest(publicKeyID, conf.PrivateKey, req, nil)
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req = req.WithContext(c)
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
	logger.Debug("getActivityJSON.response", zap.Int("code", resp.StatusCode), zap.String("body", string(body)))
	return body, nil
}

func postActivityJSON(ctx context.Context, actor Actor, url string, body []byte) ([]byte, error) {
	conf := config.FromContext(ctx)
	logger := logging.FromContext(ctx)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", mimeTypeActivityJSON)
	publicKeyID := GetPublicKeyID(actor)
	SignRequest(publicKeyID, conf.PrivateKey, req, body)
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
	logger.Debug("postActivityJSON.response", zap.Int("code", resp.StatusCode), zap.String("body", string(respBody)))
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusCreated:
	default:
		return nil, fmt.Errorf("http post status: %d", resp.StatusCode)
	}
	return respBody, nil
}
