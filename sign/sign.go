package sign

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/go-fed/httpsig"
)

const (
	DefaultPublicKeyID = "main-key"
)

func NewHeaderSigner() (httpsig.Signer, error) {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA512, httpsig.RSA_SHA256}
	headersToSign := []string{httpsig.RequestTarget, "date", "digest"}
	signer, _, err := httpsig.NewSigner(prefs, httpsig.DigestSha256, headersToSign, httpsig.Signature, 0)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func ExportPublicKey(x rsa.PublicKey) string {
	y, err := x509.MarshalPKIXPublicKey(x)
	if err != nil {
		log.Fatal(err)
	}
	z := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: y,
		},
	)
	return string(z)
}
