package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

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

func SignHeaders(payload []byte, inbox string, privateKey *rsa.PrivateKey, publicKeyID string) (map[string]string, error) {
	u, err := url.Parse(inbox)
	if err != nil {
		return nil, err
	}
	strTime := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	y := sha256.Sum256(payload)
	digest := base64.StdEncoding.EncodeToString(y[:])
	z := sha256.Sum256(
		[]byte(
			"(request-target): post " + u.Path + "\n" +
				"host: " + u.Hostname() + "\n" +
				"date: " + strTime + "\n" +
				"digest: SHA-256=" + digest,
		),
	)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, z[:])
	if err != nil {
		return nil, err
	}
	b64 := base64.StdEncoding.EncodeToString(sig[:])
	headers := map[string]string{
		"Host":   u.Hostname(),
		"Date":   strTime,
		"Digest": "SHA-256=" + digest,
		"Signature": "keyId=\"" + publicKeyID + "\"," +
			"algorithm=\"rsa-sha256\"," +
			"headers=\"(request-target) host date digest\"," +
			"signature=\"" + b64 + "\"",
		"Accept":          "application/activity+json",
		"Content-Type":    "application/activity+json",
		"Accept-Encoding": "gzip",
	}
	return headers, nil
}

func ImportPrivateKey(s string) (*rsa.PrivateKey, error) {
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
