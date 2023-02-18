package ap

var (
	ContextURIs = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		personPropertiesNS,
	}
	personPropertiesNS = map[string]string{
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
	}
)
