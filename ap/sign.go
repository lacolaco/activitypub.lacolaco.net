package ap

import (
	"crypto/rsa"
	"net/http"

	"humungus.tedunangst.com/r/webs/httpsig"
)

const (
	DefaultPublicKeyID = "key"
)

func SignRequest(publicKeyID string, key *rsa.PrivateKey, req *http.Request, body []byte) {
	httpsig.SignRequest(publicKeyID, httpsig.PrivateKey{Key: key, Type: httpsig.RSA}, req, body)
}
