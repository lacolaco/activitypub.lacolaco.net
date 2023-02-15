package usecase

import (
	"context"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/tracing"
	"github.com/lacolaco/activitypub.lacolaco.net/webfinger"
)

type searchUsecase struct {
}

func NewSearchUsecase() *searchUsecase {
	return &searchUsecase{}
}

func (u *searchUsecase) SearchPerson(ctx context.Context, id string) (*goap.Person, error) {
	ctx, span := tracing.StartSpan(ctx, "usecase.search.SearchPerson")
	defer span.End()

	personURI, err := webfinger.ResolveAccountURI(ctx, id)
	if err != nil {
		return nil, err
	}
	if personURI == "" {
		return nil, nil
	}
	person, err := ap.GetPerson(ctx, personURI)
	if err != nil {
		return nil, err
	}
	return person, nil
}
