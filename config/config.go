package config

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"golang.org/x/oauth2/google"
	"humungus.tedunangst.com/r/webs/httpsig"
)

type Config struct {
	Port              string
	PrivateKey        *httpsig.PrivateKey
	PublicKey         *httpsig.PublicKey
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
	privateKey, err := sign.DecodePrivateKey(rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public().(*rsa.PublicKey)
	config.PrivateKey = &httpsig.PrivateKey{Type: httpsig.RSA, Key: privateKey}
	config.PublicKey = &httpsig.PublicKey{Type: httpsig.RSA, Key: publicKey}

	config.googleCredentials = findGoogleCredentials()
	config.isRunningOnCloud = os.Getenv("K_SERVICE") != ""
	return &config, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}
