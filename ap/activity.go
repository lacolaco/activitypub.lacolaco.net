package ap

import (
	"context"
	"fmt"
	"net/url"
	"time"

	goap "github.com/go-ap/activitypub"
	"humungus.tedunangst.com/r/webs/junk"
)

func Accept(ctx context.Context, actor Actor, req *goap.Activity) error {
	now := time.Now()
	userID := actor.GetID()
	to, err := GetPerson(ctx, req.Actor.GetID().String())
	if err != nil {
		return err
	}

	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = fmt.Sprintf("%s/%d", userID, now.Unix())
	j["type"] = "Accept"
	j["actor"] = actor.GetID()
	j["to"] = to.GetID().String()
	j["published"] = now.UTC().Format(time.RFC3339)
	j["object"] = req

	if _, err := signedPost(ctx, actor, string(to.Inbox.GetLink()), j.ToBytes()); err != nil {
		return err
	}
	return nil
}

func GetPerson(ctx context.Context, id string) (*goap.Person, error) {
	addr, _ := url.Parse(id)
	if addr.Scheme == "" {
		addr.Scheme = "https"
	}
	body, err := signedGet(ctx, systemActor, addr.String())
	if err != nil {
		return nil, err
	}
	item, err := goap.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	var p *goap.Person
	err = goap.OnActor(item, func(a *goap.Actor) error {
		p = a
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newFollowJunk(actor Actor, whom *goap.Person) junk.Junk {
	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = fmt.Sprintf("%s/follow/%s", actor.GetID(), whom.GetID())
	j["type"] = "Follow"
	j["actor"] = actor.GetID()
	j["object"] = whom.GetID().String()
	return j
}

func FollowPerson(ctx context.Context, actor Actor, whom *goap.Person) error {
	j := newFollowJunk(actor, whom)

	if _, err := signedPost(ctx, actor, string(whom.Inbox.GetLink()), j.ToBytes()); err != nil {
		return err
	}
	return nil
}

func UnfollowPerson(ctx context.Context, actor Actor, whom *goap.Person) error {
	now := time.Now()
	undoID := fmt.Sprintf("%s/follow/%s/undo/%d", actor.GetID(), whom.GetID(), now.Unix())

	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = undoID
	j["type"] = "Undo"
	j["actor"] = actor.GetID()
	j["object"] = newFollowJunk(actor, whom)

	if _, err := signedPost(ctx, actor, string(whom.Inbox.GetLink()), j.ToBytes()); err != nil {
		return err
	}
	return nil
}
