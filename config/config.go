package config

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port              string
	RsaPrivateKey     *rsa.PrivateKey
	googleCredentials *google.Credentials
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
	rsaPrivateKey := os.Getenv("RSA_PRIVATE_KEY")
	if rsaPrivateKey == "" {
		return nil, fmt.Errorf("RSA keys are not set")
	}
	privateKey, err := sign.ImportPrivateKey(rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	config.RsaPrivateKey = privateKey

	config.googleCredentials = findGoogleCredentials()
	config.isRunningOnCloud = os.Getenv("K_SERVICE") != ""
	return &config, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}
