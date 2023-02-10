package config

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
)

type Config struct {
	Port    string
	RsaKeys struct {
		PrivateKey string
		PublicKey  string
	}

	googleCredentials *google.Credentials
}

func (c *Config) ProjectID() string {
	if c.googleCredentials != nil {
		return c.googleCredentials.ProjectID
	}
	return ""
}

func Load() (*Config, error) {
	var config Config
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	config.RsaKeys.PrivateKey = os.Getenv("RSA_PRIVATE_KEY")
	config.RsaKeys.PublicKey = os.Getenv("RSA_PUBLIC_KEY")
	if config.RsaKeys.PrivateKey == "" || config.RsaKeys.PublicKey == "" {
		return nil, fmt.Errorf("RSA keys are not set")
	}
	config.googleCredentials = findGoogleCredentials()
	return &config, nil
}

func findGoogleCredentials() *google.Credentials {
	cred, _ := google.FindDefaultCredentials(context.Background())
	return cred
}
