package ap

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"humungus.tedunangst.com/r/webs/junk"
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

func Accept(ctx context.Context, actor *Person, object *Activity) error {
	now := time.Now()
	to, err := GetPerson(ctx, object.Actor)
	if err != nil {
		return err
	}

	j := junk.New()
	j["@context"] = ContextURIs
	j["id"] = fmt.Sprintf("%s/%d", actor.ID, now.Unix())
	j["type"] = "Accept"
	j["actor"] = actor.ID
	j["to"] = to.GetID()
	j["published"] = now.UTC().Format(time.RFC3339)
	j["object"] = object

	if _, err := signedPost(ctx, actor.PublicKey, to.Inbox, j.ToBytes()); err != nil {
		return err
	}
	return nil
}

func GetPerson(ctx context.Context, id string) (*Person, error) {
	addr, _ := url.Parse(id)
	if addr.Scheme == "" {
		addr.Scheme = "https"
	}
	body, err := signedGet(ctx, systemActor.PublicKey, addr.String())
	if err != nil {
		return nil, err
	}
	p, err := PersonFromBytes(body)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newFollowJunk(from *Person, to *Person) junk.Junk {
	j := junk.New()
	j["@context"] = ContextURIs
	j["id"] = fmt.Sprintf("%s/follow/%s", from.GetID(), to.GetID())
	j["type"] = "Follow"
	j["actor"] = from.GetID()
	j["object"] = to.GetID()
	return j
}

func FollowPerson(ctx context.Context, from, to *Person) error {
	j := newFollowJunk(from, to)

	if _, err := signedPost(ctx, from.PublicKey, to.Inbox, j.ToBytes()); err != nil {
		return err
	}
	return nil
}

func UnfollowPerson(ctx context.Context, from, to *Person) error {
	now := time.Now()
	undoID := fmt.Sprintf("%s/follow/%s/undo/%d", from.GetID(), to.GetID(), now.Unix())

	j := junk.New()
	j["@context"] = ContextURIs
	j["id"] = undoID
	j["type"] = "Undo"
	j["actor"] = from.GetID()
	j["object"] = newFollowJunk(from, to)

	if _, err := signedPost(ctx, from.PublicKey, to.Inbox, j.ToBytes()); err != nil {
		return err
	}
	return nil
}
