package ap

type ObjectType string

const (
	ObjectTypeIRI ObjectType = "IRI"
)

type ActivityStreamsObject interface {
	GetType() ObjectType
}

type IRI string

var _ LinkOrIRI = IRI("")
var _ ObjectOrLink = IRI("")

func (i IRI) GetID() IRI {
	return i
}

func (i IRI) GetType() ObjectType {
	return ObjectTypeIRI
}

func (i IRI) GetLink() IRI {
	return i
}

func (i IRI) IsLink() bool {
	return true
}

func (i IRI) IsObject() bool {
	return true
}

type Object interface {
	ActivityStreamsObject
	GetID() IRI
}

type Link interface {
	GetHref() IRI
	GetName() string
}

type LinkOrIRI interface {
	GetLink() IRI
}

type ObjectOrLink interface {
	Object
	LinkOrIRI
	IsLink() bool
	IsObject() bool
}

type ActivityObject interface {
	Object
	GetActor() ObjectOrLink
	GetObject() Object
}
