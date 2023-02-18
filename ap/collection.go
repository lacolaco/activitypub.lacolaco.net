package ap

import (
	"encoding/json"

	"humungus.tedunangst.com/r/webs/junk"
)

const (
	ObjectTypeOrderedCollection ObjectType = "OrderedCollection"
)

type Item = ObjectOrLink

type OrderedCollection struct {
	Context      interface{} `json:"@context,omitempty"`
	ID           IRI         `json:"id,omitempty"`
	Type         ObjectType  `json:"type,omitempty"`
	TotalItems   int         `json:"totalItems,omitempty"`
	First        IRI         `json:"first,omitempty"`
	Last         IRI         `json:"last,omitempty"`
	OrderedItems []Item      `json:"orderedItems,omitempty"`
}

var _ ObjectOrLink = (*OrderedCollection)(nil)

func (o *OrderedCollection) GetID() IRI {
	return o.ID
}

func (o *OrderedCollection) GetType() ObjectType {
	return ObjectTypeOrderedCollection
}

func (o *OrderedCollection) IsLink() bool {
	return false
}

func (o *OrderedCollection) IsObject() bool {
	return true
}

func (o *OrderedCollection) GetLink() IRI {
	return o.GetID()
}

func (o *OrderedCollection) MarshalJSON() ([]byte, error) {
	type temp OrderedCollection
	v := temp(*o)
	v.Type = o.GetType()
	return json.Marshal(v)
}

func (o *OrderedCollection) UnmarshalJSON(data []byte) error {
	m, err := junk.FromBytes(data)
	if err != nil {
		return err
	}
	if id, ok := m.GetString("id"); ok {
		o.ID = IRI(id)
	}
	if totalItems, ok := m.GetNumber("totalItems"); ok {
		o.TotalItems = int(totalItems)
	}
	if first, ok := m.GetString("first"); ok {
		o.First = IRI(first)
	}
	if last, ok := m.GetString("last"); ok {
		o.Last = IRI(last)
	}
	if orderedItems, ok := m.GetArray("orderedItems"); ok {
		o.OrderedItems = make([]Item, 0, len(orderedItems))
		for _, item := range orderedItems {
			if iri, ok := item.(string); ok {
				o.OrderedItems = append(o.OrderedItems, IRI(iri))
				continue
			}
		}
	}
	return nil
}
