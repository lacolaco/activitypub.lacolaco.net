package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port    string
	RsaKeys struct {
		PrivateKey string
		PublicKey  string
	}
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
	return &config, nil
}
