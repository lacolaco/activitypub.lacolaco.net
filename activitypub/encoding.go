package activitypub

import (
	goap "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

var (
	contextURIs = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		map[string]interface{}{
			"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
			"sensitive":                 "as:sensitive",
			"Hashtag":                   "as:Hashtag",
			"quoteUrl":                  "as:quoteUrl",
			"toot":                      "http://joinmastodon.org/ns#",
			"discoverable":              "toot:discoverable",
			"Emoji":                     "toot:Emoji",
			"featured":                  "toot:featured",
			"misskey":                   "https://misskey-hub.net/ns#",
			"schema":                    "http://schema.org#",
			"PropertyValue":             "schema:PropertyValue",
			"value":                     "schema:value",
		},
	}
)

func MarshalActivityJSON(v interface{}) ([]byte, error) {
	return jsonld.WithContext(jsonld.IRI(goap.ActivityBaseURI), jsonld.IRI(goap.SecurityContextURI)).Marshal(v)
}
