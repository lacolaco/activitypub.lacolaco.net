package ap

import (
	"context"
	"fmt"
	"time"

	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"humungus.tedunangst.com/r/webs/junk"
)

func Accept(ctx context.Context, actor Actor, req *goap.Activity) error {
	now := time.Now()
	userID := actor.GetID()
	to, err := getPerson(ctx, req.Actor.GetID().String())
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

	if _, err := postActivityJSON(ctx, actor, string(to.Inbox.GetLink()), j.ToBytes()); err != nil {
		return err
	}
	return nil
}

func GetPerson(ctx context.Context, id string) (*goap.Person, error) {
	return getPerson(ctx, id)
}

func FollowPerson(ctx context.Context, actor Actor, target string) (*model.Job, error) {
	now := time.Now()
	to, err := getPerson(ctx, target)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s/%d", actor.GetID(), now.Unix())

	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = id
	j["type"] = "Follow"
	j["actor"] = actor.GetID()
	j["object"] = to.GetID().String()

	if _, err := postActivityJSON(ctx, actor, string(to.Inbox.GetLink()), j.ToBytes()); err != nil {
		return nil, err
	}
	return model.NewJob(id, model.JobTypeFollowUser, actor.GetID(), target), nil
}

func UnfollowPerson(ctx context.Context, actor Actor, target string) (*model.Job, error) {
	now := time.Now()
	to, err := getPerson(ctx, target)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s/%d", actor.GetID(), now.Unix())

	obj := junk.New()
	obj["@context"] = contextURIs
	obj["id"] = id
	obj["type"] = "Follow"
	obj["actor"] = actor.GetID()
	obj["object"] = to.GetID().String()

	j := junk.New()
	j["@context"] = contextURIs
	j["id"] = id
	j["type"] = "Undo"
	j["actor"] = actor.GetID()
	j["object"] = obj

	if _, err := postActivityJSON(ctx, actor, string(to.Inbox.GetLink()), j.ToBytes()); err != nil {
		return nil, err
	}
	return model.NewJob(id, model.JobTypeUnfollowUser, actor.GetID(), target), nil
}
