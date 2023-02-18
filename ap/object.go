package ap

import "encoding/json"

const ObjectTypeImage ObjectType = "Image"

// Image is an image
// See https://www.w3.org/TR/activitystreams-vocabulary/#dfn-image
type Image struct {
	Type      ObjectType `json:"type,omitempty"`
	URL       string     `json:"url,omitempty"`
	MediaType string     `json:"mediaType,omitempty"`
}

var _ LinkOrIRI = (*Image)(nil)

func (i *Image) GetType() ObjectType {
	return ObjectTypeImage
}

func (i *Image) GetLink() IRI {
	return IRI(i.URL)
}

func (i *Image) MarshalJSON() ([]byte, error) {
	type temp Image
	v := temp(*i)
	v.Type = i.GetType()
	return json.Marshal(v)
}

const ObjectTypePropertyValue ObjectType = "PropertyValue"

// PropertyValue is a property-value pair
// See https://schema.org/PropertyValue
type PropertyValue struct {
	Type  ObjectType `json:"type,omitempty"`
	Name  string     `json:"name,omitempty"`
	Value string     `json:"value,omitempty"`
}

var _ ActivityStreamsObject = (*PropertyValue)(nil)

func (p *PropertyValue) GetType() ObjectType {
	return ObjectTypePropertyValue
}

func (p *PropertyValue) MarshalJSON() ([]byte, error) {
	type temp PropertyValue
	v := temp(*p)
	v.Type = p.GetType()
	return json.Marshal(v)
}
