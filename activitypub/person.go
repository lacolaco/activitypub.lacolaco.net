package activitypub

import (
	"fmt"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/sign"
	"humungus.tedunangst.com/r/webs/httpsig"
)

func NewPersonJSON(u *model.User, baseUri string, publicKey *httpsig.PublicKey) map[string]interface{} {
	apID := u.GetActivityPubID(baseUri)
	publicKeyPem, err := httpsig.EncodeKey(publicKey.Key)
	if err != nil {
		panic(err)
	}

	return map[string]interface{}{
		"@context":                  contextURIs,
		"type":                      "Person",
		"id":                        apID,
		"name":                      u.Name,
		"preferredUsername":         u.PrefName,
		"summary":                   u.Description,
		"inbox":                     fmt.Sprintf("%s/inbox", apID),
		"outbox":                    fmt.Sprintf("%s/outbox", apID),
		"followers":                 fmt.Sprintf("%s/followers", apID),
		"following":                 fmt.Sprintf("%s/following", apID),
		"url":                       fmt.Sprintf("%s/@%s", baseUri, u.ID),
		"published":                 u.CreatedAt.Format(time.RFC3339),
		"updated":                   u.UpdatedAt.Format(time.RFC3339),
		"discoverable":              true,
		"manuallyApprovesFollowers": false,
		"icon": map[string]interface{}{
			"type":      "Image",
			"mediaType": u.Icon.MediaType,
			"url":       u.Icon.URL,
		},
		"publicKey": map[string]interface{}{
			"id":           fmt.Sprintf("%s#%s", apID, sign.DefaultPublicKeyID),
			"owner":        apID,
			"publicKeyPem": publicKeyPem,
		},
	}
}
