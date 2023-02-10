package activitypub

import (
	goap "github.com/go-ap/activitypub"
)

func NewAccept(id string, actor string, target goap.Item) *goap.Accept {
	return &goap.Accept{
		Context: goap.ActivityBaseURI,
		ID:      goap.IRI(id),
		Type:    goap.AcceptType,
		Actor:   goap.IRI(actor),
		Object:  target,
	}
}
