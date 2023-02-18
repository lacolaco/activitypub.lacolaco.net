package ap_test

import (
	"encoding/json"
	"testing"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"humungus.tedunangst.com/r/webs/junk"
)

func TestActivity(t *testing.T) {
	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		a := &ap.Activity{
			ID:     ap.IRI("https://example.com/activity"),
			Type:   ap.ActivityTypeFollow,
			Actor:  ap.IRI("https://example.com/actor"),
			Object: ap.IRI("https://example.com/object"),
		}
		b, err := a.MarshalJSON()
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
	})

	t.Run("can be unmarshalled from JSON", func(tt *testing.T) {
		b := []byte(`{"id": "https://example.com/activity", "type": "Follow", "actor": "https://example.com/actor", "object": "https://example.com/object"}`)
		var a ap.Activity
		if err := json.Unmarshal(b, &a); err != nil {
			tt.Fatal(err)
		}
		if string(a.ID) != "https://example.com/activity" {
			tt.Error("id is not https://example.com/activity", a.ID)
		}
		if a.Type != ap.ActivityTypeFollow {
			tt.Error("type is not Follow", a.Type)
		}
		if string(a.Actor.GetID()) != "https://example.com/actor" {
			tt.Error("actor is not https://example.com/actor", a.Actor)
		}
		if string(a.Object.GetID()) != "https://example.com/object" {
			tt.Error("object is not https://example.com/object", a.Object)
		}
	})

	t.Run("can be unmarshalled from unsupported type JSON", func(tt *testing.T) {
		b := []byte(`{"id": "https://example.com/activity", "type": "Note", "actor": "https://example.com/actor", "object": "https://example.com/object"}`)
		var a ap.Activity
		if err := json.Unmarshal(b, &a); err != nil {
			tt.Fatal(err)
		}
		if string(a.ID) != "https://example.com/activity" {
			tt.Error("id is not https://example.com/activity", a.ID)
		}
		if a.Type != "Note" {
			tt.Error("type is not Follow", a.Type)
		}
		if string(a.Actor.GetID()) != "https://example.com/actor" {
			tt.Error("actor is not https://example.com/actor", a.Actor)
		}
		if string(a.Object.GetID()) != "https://example.com/object" {
			tt.Error("object is not https://example.com/object", a.Object)
		}
	})
}
