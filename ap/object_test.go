package ap_test

import (
	"encoding/json"
	"testing"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"humungus.tedunangst.com/r/webs/junk"
)

func TestImage(t *testing.T) {
}

func TestPropertyValue(t *testing.T) {

	t.Run("can be marshalled to JSON", func(tt *testing.T) {
		p := &ap.PropertyValue{
			Name:  "name",
			Value: "value",
		}
		b, err := json.Marshal(p)
		if err != nil {
			tt.Fatal(err)
		}
		j, err := junk.FromBytes(b)
		if err != nil {
			tt.Fatal(err)
		}
		if dType, ok := j["type"]; !ok || dType != "PropertyValue" {
			tt.Error("type is not included or not PropertyValue", dType)
		}
		if _, ok := j["name"]; !ok {
			tt.Error("name is not included")
		}
		if _, ok := j["value"]; !ok {
			tt.Error("value is not included")
		}
	})

	t.Run("can be unmarshalled from JSON", func(tt *testing.T) {
		b := []byte(`{"type": "PropertyValue", "name": "name", "value": "value"}`)
		p := &ap.PropertyValue{}
		if err := json.Unmarshal(b, p); err != nil {
			tt.Fatal(err)
		}
		if p.Name != "name" {
			tt.Error("name is not name", p.Name)
		}
		if p.Value != "value" {
			tt.Error("value is not value", p.Value)
		}
	})
}
