package ap

import (
	"crypto/rsa"
	"net/http"

	"humungus.tedunangst.com/r/webs/httpsig"
)

func SignRequest(publicKeyID string, key *rsa.PrivateKey, req *http.Request, body []byte) {
	httpsig.SignRequest(publicKeyID, httpsig.PrivateKey{Key: key, Type: httpsig.RSA}, req, body)
}

func VerifyRequest(req *http.Request, content []byte) error {
	if _, err := httpsig.VerifyRequest(req, content, httpsig.ActivityPubKeyGetter); err != nil {
		return err
	}
	return nil
}
