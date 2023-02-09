package activitypub_test

import (
	"encoding/json"
	"testing"

	"github.com/lacolaco/activitypub.lacolaco.net/activitypub"
)

func TestActor(t *testing.T) {

	t.Run("unmarshal json", func(tt *testing.T) {
		body := `{
			"@context":["https://www.w3.org/ns/activitystreams","https://w3id.org/security/v1",{
				"manuallyApprovesFollowers":"as:manuallyApprovesFollowers","toot":"http://joinmastodon.org/ns#",
				"featured":{"@id":"toot:featured","@type":"@id"},
				"alsoKnownAs":{"@id":"as:alsoKnownAs","@type":"@id"},
				"movedTo":{"@id":"as:movedTo","@type":"@id"},
				"schema":"http://schema.org#","PropertyValue":"schema:PropertyValue","value":"schema:value",
				"IdentityProof":"toot:IdentityProof","discoverable":"toot:discoverable","focalPoint":{"@container":"@list","@id":"toot:focalPoint"}}],
			"id":"https://activitypub.lacolaco.net/users/lacolaco",
			"type":"Person",
			"following":"https://activitypub.lacolaco.net/users/lacolaco/following",
			"followers":"https://activitypub.lacolaco.net/users/lacolaco/followers",
			"inbox":"https://activitypub.lacolaco.net/users/lacolaco/inbox",
			"outbox":"https://activitypub.lacolaco.net/users/lacolaco/outbox",
			"featured":"https://activitypub.lacolaco.net/users/lacolaco/collections/featured",
			"preferredUsername":"lacolaco",
			"name":"らこ",
			"summary":"\u003cp\u003eらこだよ。\u003c/p\u003e",
			"url":"https://activitypub.lacolaco.net/@lacolaco",
			"manuallyApprovesFollowers":false,
			"discoverable":null,
			"publicKey":{
				"id":"https://activitypub.lacolaco.net/users/lacolaco#main-key",
				"owner":"https://activitypub.lacolaco.net/users/lacolaco",
				"publicKeyPem":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmfote2wqxEcqD3n3IPXs\nuv8ngs0hQ1fzvydoAQLntwW8ch9bJLJhjT4Z2VPVqxmrfM0l4keSaf7AesevE/CK\n1trfK3zFQvumzU1J8Gu04+1EyI0wrv5d1SpVHR+xrvxY4Z4UsKVb/wrcKLITV7SW\nZzzhUMSoyAJuwc0V+s8R0hBYEpleq/02mrSvyeC2QoTj7+vqEO5KIMipBn9hRFm/\nA0ba78zbjScWFj6G8RqK2NGJ/dZk9S6S/zqIvQXwSp+Cyi5EepsLUqqGnvb0Wyi9\nFqZD2PVqksYZM1bNJAh8vN8rkwbi4/wNXzptHnSIJ0c4GEado0kcnY+Mbp9LTbPn\nGwIDAQAB\n-----END PUBLIC KEY-----\n"},
			"tag":[],
			"attachment":[
				{"type":"PropertyValue","name":"Twitter",
				"value":"\u003ca href=\"https://twitter.com/laco2net\" rel=\"me nofollow noopener\" target=\"_blank\"\u003e\u003cspan class=\"invisible\"\u003ehttps://\u003c/span\u003e\u003cspan class=\"\"\u003etwitter.com/laco2net\u003c/span\u003e\u003cspan class=\"invisible\"\u003e\u003c/span\u003e\u003c/a\u003e"},{"type":"PropertyValue","name":"Website","value":"\u003ca href=\"https://lacolaco.net\" rel=\"me nofollow noopener\" target=\"_blank\"\u003e\u003cspan class=\"invisible\"\u003ehttps://\u003c/span\u003e\u003cspan class=\"\"\u003elacolaco.net\u003c/span\u003e\u003cspan class=\"invisible\"\u003e\u003c/span\u003e\u003c/a\u003e"}
			],
			"endpoints":{
				"sharedInbox":"https://activitypub.lacolaco.net/inbox"
			},
			"icon":{
				"type":"Image",
				"mediaType":"image/png",
				"url":"https://storage.googleapis.com/lacolaco-mastodon/accounts/avatars/000/000/001/original/c475c32751b3d022.png"
			}
		}
		`
		actor := &activitypub.Actor{}
		err := json.Unmarshal([]byte(body), actor)
		if err != nil {
			tt.Fatal(err)
		}
		if actor.ID != "https://activitypub.lacolaco.net/users/lacolaco" {
			tt.Fatal("ID is not matched")
		}
		if actor.Type != "Person" {
			tt.Fatal("Type is not matched")
		}
	})
}
