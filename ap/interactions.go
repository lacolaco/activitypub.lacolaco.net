package ap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var (
	systemActor = &Person{
		ID: "https://activitypub.lacolaco.net",
		PublicKey: &PublicKey{
			ID:    "https://activitypub.lacolaco.net#main-key",
			Owner: "https://activitypub.lacolaco.net",
		},
	}
)

func Accept(ctx context.Context, actor *Person, object ActivityObject) error {
	now := time.Now()
	to, err := GetPerson(ctx, object.GetActor().GetID())
	if err != nil {
		return err
	}

	accept := NewActivityAccept()
	accept.Context = ContextURIs
	accept.ID = IRI(fmt.Sprintf("%s/accept/%s", actor.ID, object.GetID()))
	accept.Actor = actor.ID
	accept.Object = object
	accept.To = []Item{to.ID}
	accept.Published = now.UTC()

	b, err := json.Marshal(accept)
	if err != nil {
		return err
	}
	if _, err := signedPost(ctx, actor.PublicKey, to.Inbox, b); err != nil {
		return err
	}
	return nil
}

func GetPerson(ctx context.Context, id IRI) (*Person, error) {
	body, err := signedGet(ctx, systemActor.PublicKey, id)
	if err != nil {
		return nil, err
	}
	var p *Person
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

func FollowPerson(ctx context.Context, from, to *Person) error {
	activity := NewActivityFollow()
	activity.Context = ContextURIs
	activity.ID = IRI(fmt.Sprintf("%s/follow/%s", from.GetID(), to.GetID()))
	activity.Actor = from.GetID()
	activity.Object = to.GetID()

	b, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	if _, err := signedPost(ctx, from.PublicKey, to.GetInbox(), b); err != nil {
		return err
	}
	return nil
}

func UnfollowPerson(ctx context.Context, from, to *Person) error {
	follow := NewActivityFollow()
	follow.ID = IRI(fmt.Sprintf("%s/follow/%s", from.GetID(), to.GetID()))
	follow.Actor = from.GetID()
	follow.Object = to.GetID()

	now := time.Now()
	activity := NewActivityUndo()
	activity.Context = ContextURIs
	activity.ID = IRI(fmt.Sprintf("%s/follow/%s/undo/%d", from.GetID(), to.GetID(), now.Unix()))
	activity.Actor = from.GetID()
	activity.Object = follow

	b, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	if _, err := signedPost(ctx, from.PublicKey, to.GetInbox(), b); err != nil {
		return err
	}
	return nil
}
