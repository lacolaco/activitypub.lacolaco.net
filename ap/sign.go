package ap

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"humungus.tedunangst.com/r/webs/httpsig"
)

const (
	publicKeyIDSuffix = "key"
)

func SignRequest(publicKeyID string, key *rsa.PrivateKey, req *http.Request, body []byte) {
	httpsig.SignRequest(publicKeyID, httpsig.PrivateKey{Key: key, Type: httpsig.RSA}, req, body)
}

func GetPublicKeyID(actor Actor) string {
	return fmt.Sprintf("%s#%s", actor.GetID(), publicKeyIDSuffix)
}
