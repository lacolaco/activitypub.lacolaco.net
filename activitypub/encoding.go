package activitypub

import (
	goap "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
)

func MarshalActivityJSON(v interface{}) ([]byte, error) {
	return jsonld.WithContext(jsonld.IRI(goap.ActivityBaseURI), jsonld.IRI(goap.SecurityContextURI)).Marshal(v)
}
