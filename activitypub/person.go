package activitypub

import (
	"fmt"

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
		"@context":          contextURIs,
		"type":              "Person",
		"id":                apID,
		"name":              u.Name,
		"preferredUsername": u.PrefName,
		"summary":           u.Description,
		"inbox":             fmt.Sprintf("%s/inbox", apID),
		"outbox":            fmt.Sprintf("%s/outbox", apID),
		"followers":         fmt.Sprintf("%s/followers", apID),
		"following":         fmt.Sprintf("%s/following", apID),
		// "featured":                  fmt.Sprintf("%s/collections/featured", apID),
		"discoverable":              true,
		"manuallyApprovesFollowers": false,
		// "sharedInbox":               fmt.Sprintf("%s/inbox", baseUri),
		// "endpoints": map[string]interface{}{
		// 	"sharedInbox": fmt.Sprintf("%s/inbox", baseUri),
		// },
		"url": fmt.Sprintf("%s/@%s", baseUri, u.ID),
		"icon": map[string]interface{}{
			"type":      "Image",
			"mediaType": u.Icon.MediaType,
			"uRL":       u.Icon.URL,
		},
		"publicKey": map[string]interface{}{
			"id":           fmt.Sprintf("%s#%s", apID, sign.DefaultPublicKeyID),
			"owner":        apID,
			"publicKeyPem": publicKeyPem,
		},
	}
}
