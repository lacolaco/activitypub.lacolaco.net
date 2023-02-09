package sign

import (
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
