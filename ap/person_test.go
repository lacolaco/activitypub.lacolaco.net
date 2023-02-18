package ap_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"humungus.tedunangst.com/r/webs/junk"
)

func TestPerson(t *testing.T) {

	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		person := &ap.Person{
			ID:                ap.IRI("https://example.com/person"),
			Name:              "Test User",
			PreferredUsername: "test",
			Summary:           "This is a test user",
			Icon: &ap.Image{
				URL:       "https://example.com/icon.png",
				MediaType: "image/png",
			},
			Published: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			PublicKey: &ap.PublicKey{
				ID:           "https://example.com/person#main-key",
				Owner:        ap.IRI("https://example.com/person"),
				PublicKeyPem: "PUBLIC KEY",
			},
			Attachment: []ap.ActivityStreamsObject{
				&ap.PropertyValue{
					Name:  "test",
					Value: "value",
				},
			},
		}
		b, err := person.MarshalJSON()
		if err != nil {
			tt.Fatal(err)
		}
		j, err := junk.FromBytes(b)
		if err != nil {
			tt.Fatal(err)
		}
		t.Logf("%s", j.ToString())
		if _, ok := j["id"]; !ok {
			tt.Error("id is not included")
		}
		if dType, ok := j["type"]; !ok || dType != "Person" {
			tt.Error("type is not included or not Person", dType)
		}
		if _, ok := j["name"]; !ok {
			tt.Error("name is not included")
		}
		// attachment
		attachment, ok := j["attachment"]
		if !ok {
			tt.Error("attachment is not included")
		}
		if _, ok := attachment.([]interface{}); !ok {
			tt.Error("attachment is not array")
		}
		// toot property
		if _, ok := j["discoverable"]; !ok {
			tt.Error("discoverable is not included")
		}
		// custom property
		if _, ok := j["manuallyApprovesFollowers"]; !ok {
			tt.Error("manuallyApprovesFollowers is not included")
		}
	})

	t.Run("can be unmarshalled from JSON", func(tt *testing.T) {
		b := []byte(`{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id": "https://example.com/users/test",
			"type": "Person",
			"name": "Test User",
			"preferredUsername": "Test",
			"summary": "This is a test user",
			"inbox": "https://example.com/users/test/inbox",
			"outbox": "https://example.com/users/test/outbox",
			"followers": "https://example.com/users/test/followers",
			"following": "https://example.com/users/test/following",
			"url": "https://example.com/@test",
			"published": "2020-01-01T00:00:00Z",
			"icon": {
				"type": "Image",
				"url": "https://example.com/icon.png",
				"mediaType": "image/png"
			},
			"publicKey": {
				"id": "https://example.com/users/test#main-key",
				"owner": "https://example.com/users/test",
				"publicKeyPem": "----"
			},
			"attachment": [
				{
					"type": "PropertyValue",
					"name": "Twitter",
					"value": "<a href=\"https://twitter.com/test\" rel=\"me\">@test</a>"
				}
			],
			"discoverable": true,
			"manuallyApprovesFollowers": true
		}`)
		var person *ap.Person
		if err := json.Unmarshal(b, &person); err != nil {
			tt.Fatal(err)
		}
		if person.ID != "https://example.com/users/test" {
			tt.Error("ID is not deserialized")
		}
		if person.Name != "Test User" {
			tt.Error("Name is not deserialized")
		}
		if person.PreferredUsername != "Test" {
			tt.Error("PreferredUsername is not deserialized")
		}
		if person.Summary != "This is a test user" {
			tt.Error("Summary is not deserialized")
		}
		if person.Inbox != "https://example.com/users/test/inbox" {
			tt.Error("Inbox is not deserialized")
		}
		if person.Outbox != "https://example.com/users/test/outbox" {
			tt.Error("Outbox is not deserialized")
		}
		if person.Followers != "https://example.com/users/test/followers" {
			tt.Error("Followers is not deserialized")
		}
		if person.Following != "https://example.com/users/test/following" {
			tt.Error("Following is not deserialized")
		}
		if person.URL != "https://example.com/@test" {
			tt.Error("URL is not deserialized")
		}
		if person.Published != time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) {
			tt.Error("Published is not deserialized")
		}
		if person.Icon.URL != "https://example.com/icon.png" {
			tt.Error("Icon is not deserialized")
		}
		if person.PublicKey.ID != "https://example.com/users/test#main-key" {
			tt.Error("PublicKey is not deserialized")
		}
		if person.Attachment[0].(*ap.PropertyValue).Name != "Twitter" {
			tt.Error("Attachment is not deserialized")
		}
		if person.Discoverable != true {
			tt.Error("Discoverable is not deserialized")
		}
		if person.ManuallyApprovesFollowers != true {
			tt.Error("ManuallyApprovesFollowers is not deserialized")
		}
	})
}
