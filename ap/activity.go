package ap

import (
	"context"
	"fmt"
	"time"

	goap "github.com/go-ap/activitypub"
	"humungus.tedunangst.com/r/webs/junk"
)

func Accept(ctx context.Context, user *Person, req *goap.Activity) error {
	now := time.Now()
	userID := user.ID()
	actor, err := getActor(ctx, user.PubkeyID(), req.Actor.GetID().String())
	if err != nil {
		return err
	}

	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = fmt.Sprintf("%s/%d", userID, now.Unix())
	j["type"] = "Accept"
	j["actor"] = userID
	j["to"] = actor.GetID().String()
	j["published"] = now.UTC().Format(time.RFC3339)
	j["object"] = req

	if _, err := postActivityJSON(ctx, userID, string(actor.Inbox.GetLink()), j.ToBytes()); err != nil {
		return err
	}
	return nil
}
