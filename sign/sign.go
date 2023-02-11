package sign

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
)

const (
	DefaultPublicKeyID = "main-key"
)

func DecodePrivateKey(s string) (*rsa.PrivateKey, error) {
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

func ExportPublicKey(x *rsa.PublicKey) string {
	y, err := x509.MarshalPKIXPublicKey(x)
	if err != nil {
		log.Fatalf("failed to marshal public key: %v", err)
		return ""
	}
	z := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: y,
		},
	)
	return string(z)
}
