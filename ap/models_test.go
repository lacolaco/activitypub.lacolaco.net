package ap_test

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"humungus.tedunangst.com/r/webs/junk"
)

func TestPerson(t *testing.T) {
	os.Clearenv()
	godotenv.Load("../.env")
	conf, _ := config.Load()

	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		user := &model.LocalUser{
			ID:          "test",
			UID:         "test",
			Name:        "Test User",
			PrefName:    "Test",
			Description: "This is a test user",
			CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Icon: &model.UserIcon{
				URL:       "https://example.com/icon.png",
				MediaType: "image/png",
			},
		}
		person := user.ToPerson("https://example.com", conf.PublicKey)
		b, err := person.ToBytes()
		if err != nil {
			tt.Fatal(err)
		}
		j, err := junk.FromBytes(b)
		if err != nil {
			tt.Fatal(err)
		}
		if _, ok := j["id"]; !ok {
			tt.Error("id is not included")
		}
		if dType, ok := j["type"]; !ok || dType != "Person" {
			tt.Error("type is not included or not Person", dType)
		}
		if _, ok := j["name"]; !ok {
			tt.Error("name is not included")
		}
		if _, ok := j["discoverable"]; !ok {
			tt.Error("discoverable is not included")
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
			"discoverable": true,
			"manuallyApprovesFollowers": true
		}`)
		person, err := ap.PersonFromBytes(b)
		if err != nil {
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
		if person.Icon.URL != "https://example.com/icon.png" {
			tt.Error("Icon is not deserialized")
		}
		if person.PublicKey.ID != "https://example.com/users/test#main-key" {
			tt.Error("PublicKey is not deserialized")
		}
		if person.Discoverable != true {
			tt.Error("Discoverable is not deserialized")
		}
		if person.ManuallyApprovesFollowers != true {
			tt.Error("ManuallyApprovesFollowers is not deserialized")
		}
	})
}

func TestOrderedCollection(t *testing.T) {
	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		o := &ap.OrderedCollection{
			ID: "https://example.com/collection",
			OrderedItems: []ap.ActivityPubObject{
				ap.IRI("https://example.com/1"),
				ap.IRI("https://example.com/2"),
			},
			TotalItems: 2,
		}
		b, err := o.ToBytes()
		if err != nil {
			tt.Fatal(err)
		}
		j, err := junk.FromBytes(b)
		if err != nil {
			tt.Fatal(err)
		}
		if _, ok := j["id"]; !ok {
			tt.Error("id is not included")
		}
		if dType, ok := j["type"]; !ok || dType != "OrderedCollection" {
			tt.Error("type is not included or not Person", dType)
		}
		if _, ok := j["totalItems"]; !ok {
			tt.Error("totalItems is not included")
		}
		if _, ok := j["orderedItems"]; !ok {
			tt.Error("orderedItems is not included")
		}
	})
}
