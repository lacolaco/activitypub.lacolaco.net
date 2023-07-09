package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
)

type Config struct {
	Port              string
	PrivateKey        *rsa.PrivateKey
	PublicKey         *rsa.PublicKey
	googleCredentials *google.Credentials
	ClientOrigin      string
	isRunningOnCloud  bool
}

func (c *Config) ProjectID() string {
	if c.googleCredentials != nil {
		return c.googleCredentials.ProjectID
	}
	return ""
}

func (c *Config) IsRunningOnCloud() bool {
	return c.isRunningOnCloud
}

func Load() (*Config, error) {
	var config Config
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	config.ClientOrigin = os.Getenv("CLIENT_ORIGIN")
	rsaPrivateKey := os.Getenv("RSA_PRIVATE_KEY")
	if rsaPrivateKey == "" {
		return nil, fmt.Errorf("RSA keys are not set")
	}
	privateKey, err := decodePrivateKey(rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public().(*rsa.PublicKey)
	config.PrivateKey = privateKey
	config.PublicKey = publicKey

	config.googleCredentials = findGoogleCredentials()
	config.isRunningOnCloud = os.Getenv("K_SERVICE") != ""
	return &config, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}

func decodePrivateKey(s string) (*rsa.PrivateKey, error) {
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	s = strings.Join(strings.Split(s, "\\n"), "\n")

	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key: %s", block.Type)
	}

	if block == nil {
		return nil, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}
