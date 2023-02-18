package ap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/logging"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var (
	HTTPClient = http.DefaultClient
)

const (
	userAgent            = "activitypub.lacolaco.net/1.0"
	mimeTypeActivityJSON = "application/activity+json"
)

func signedGet(ctx context.Context, publicKey *PublicKey, iri IRI) ([]byte, error) {
	ctx, span := tracing.StartSpan(ctx, "ap.signedGet")
	defer span.End()
	url, err := url.Parse(string(iri))
	if err != nil {
		return nil, err
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	span.SetAttributes(attribute.String("url", url.String()))

	conf := config.ConfigFromContext(ctx)
	logger := logging.LoggerFromContext(ctx)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mimeTypeActivityJSON)
	req.Header.Set("User-Agent", userAgent)
	SignRequest(string(publicKey.ID), conf.PrivateKey, req, nil)
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
	logger.Debug("signedGet.response", zap.Int("code", resp.StatusCode), zap.String("body", string(body)))
	return body, nil
}

func signedPost(ctx context.Context, publicKey *PublicKey, iri IRI, body []byte) ([]byte, error) {
	ctx, span := tracing.StartSpan(ctx, "ap.signedPost")
	defer span.End()

	url, err := url.Parse(string(iri))
	if err != nil {
		return nil, err
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	span.SetAttributes(attribute.String("url", url.String()))

	conf := config.ConfigFromContext(ctx)
	logger := logging.LoggerFromContext(ctx)
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", mimeTypeActivityJSON)
	SignRequest(string(publicKey.ID), conf.PrivateKey, req, body)
	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req = req.WithContext(c)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("signedPost.response", zap.Int("code", resp.StatusCode), zap.String("body", string(respBody)))
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusCreated:
	default:
		return nil, fmt.Errorf("http post status: %d", resp.StatusCode)
	}
	return respBody, nil
}
