package ap_test

import (
	"encoding/json"
	"testing"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"humungus.tedunangst.com/r/webs/junk"
)

func TestOrderedCollection(t *testing.T) {
	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		o := &ap.OrderedCollection{
			ID: "https://example.com/collection",
			OrderedItems: []ap.ObjectOrLink{
				ap.IRI("https://example.com/1"),
				ap.IRI("https://example.com/2"),
			},
			TotalItems: 2,
		}
		b, err := json.Marshal(o)
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
			tt.Error("type is not included or not OrderedCollection", dType)
		}
		if _, ok := j["totalItems"]; !ok {
			tt.Error("totalItems is not included")
		}
		if _, ok := j["orderedItems"]; !ok {
			tt.Error("orderedItems is not included")
		}
	})

	t.Run("can be unmarshalled from JSON", func(tt *testing.T) {
		b := []byte(`{"id": "https://example.com/collection", "type": "OrderedCollection", "totalItems": 2, "orderedItems": ["https://example.com/1", "https://example.com/2"]}`)
		var o ap.OrderedCollection
		if err := json.Unmarshal(b, &o); err != nil {
			tt.Fatal(err)
		}
		if o.ID != "https://example.com/collection" {
			tt.Error("id is not https://example.com/collection", o.ID)
		}
		if o.TotalItems != 2 {
			tt.Error("totalItems is not 2", o.TotalItems)
		}
		if len(o.OrderedItems) != 2 {
			tt.Error("orderedItems is not 2", len(o.OrderedItems))
		}
		if o.OrderedItems[0] != ap.IRI("https://example.com/1") {
			tt.Error("orderedItems[0] is not https://example.com/1", o.OrderedItems[0])
		}
		if o.OrderedItems[1] != ap.IRI("https://example.com/2") {
			tt.Error("orderedItems[1] is not https://example.com/2", o.OrderedItems[1])
		}
	})
}
