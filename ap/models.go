package ap

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-fed/activity/streams"
	typeactivity "github.com/go-fed/activity/streams/impl/activitystreams/type_activity"
	typeorderedcollection "github.com/go-fed/activity/streams/impl/activitystreams/type_orderedcollection"
	typeperson "github.com/go-fed/activity/streams/impl/activitystreams/type_person"
	"github.com/go-fed/activity/streams/vocab"
	"humungus.tedunangst.com/r/webs/junk"
)

type Serializable interface {
	Serialize() map[string]interface{}
}

type ActivityPubObject interface {
	GetID() string
	GetType() string
}

type IRI string

func (i IRI) GetID() string {
	return string(i)
}

func (i IRI) GetType() string {
	return "IRI"
}

type Image struct {
	URL       string `json:"url,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

type PublicKey struct {
	ID           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Person struct {
	ID                        string     `json:"id,omitempty"`
	Name                      string     `json:"name,omitempty"`
	PreferredUsername         string     `json:"preferredUsername,omitempty"`
	Summary                   string     `json:"summary,omitempty"`
	Inbox                     string     `json:"inbox,omitempty"`
	Outbox                    string     `json:"outbox,omitempty"`
	Followers                 string     `json:"followers,omitempty"`
	Following                 string     `json:"following,omitempty"`
	URL                       string     `json:"url,omitempty"`
	Published                 time.Time  `json:"published,omitempty"`
	Icon                      *Image     `json:"icon,omitempty"`
	PublicKey                 *PublicKey `json:"publicKey,omitempty"`
	Discoverable              bool       `json:"discoverable,omitempty"`
	ManuallyApprovesFollowers bool       `json:"manuallyApprovesFollowers,omitempty"`
}

func (p *Person) GetID() string {
	return p.ID
}

func (p *Person) GetType() string {
	return "Person"
}

func (p *Person) ToBytes() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	j, err := junk.FromBytes(b)
	if err != nil {
		return nil, err
	}
	j["type"] = p.GetType()
	deserialized, err := typeperson.DeserializePerson(j, map[string]string{})
	if err != nil {
		return nil, err
	}
	m, err := streams.Serialize(deserialized)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func PersonFromBytes(b []byte) (*Person, error) {
	var p *Person
	resolver, err := streams.NewJSONResolver(func(ctx context.Context, asPerson vocab.ActivityStreamsPerson) error {
		p = &Person{}
		p.ID = asPerson.GetJSONLDId().Get().String()
		p.Name = asPerson.GetActivityStreamsName().Begin().GetXMLSchemaString()
		p.PreferredUsername = asPerson.GetActivityStreamsPreferredUsername().GetXMLSchemaString()
		p.Summary = asPerson.GetActivityStreamsSummary().Begin().GetXMLSchemaString()
		p.Inbox = asPerson.GetActivityStreamsInbox().GetIRI().String()
		p.Outbox = asPerson.GetActivityStreamsOutbox().GetIRI().String()
		p.Followers = asPerson.GetActivityStreamsFollowers().GetIRI().String()
		p.Following = asPerson.GetActivityStreamsFollowing().GetIRI().String()
		p.URL = asPerson.GetActivityStreamsUrl().Begin().GetIRI().String()
		p.Published = asPerson.GetActivityStreamsPublished().Get()
		if icon := asPerson.GetActivityStreamsIcon(); !icon.Empty() {
			p.Icon = &Image{}
			switch {
			case icon.Begin().IsActivityStreamsImage():
				p.Icon.URL = icon.Begin().GetActivityStreamsImage().GetActivityStreamsUrl().Begin().GetIRI().String()
				p.Icon.MediaType = icon.Begin().GetActivityStreamsImage().GetActivityStreamsMediaType().Get()
			case icon.Begin().IsActivityStreamsLink():
				p.Icon.URL = icon.Begin().GetActivityStreamsLink().GetActivityStreamsHref().GetIRI().String()
				p.Icon.MediaType = icon.Begin().GetActivityStreamsLink().GetActivityStreamsMediaType().Get()
			}
		}
		if publicKey := asPerson.GetW3IDSecurityV1PublicKey(); !publicKey.Empty() {
			p.PublicKey = &PublicKey{}
			p.PublicKey.ID = publicKey.Begin().Get().GetJSONLDId().Get().String()
			p.PublicKey.Owner = publicKey.Begin().Get().GetW3IDSecurityV1Owner().Get().String()
			p.PublicKey.PublicKeyPem = publicKey.Begin().Get().GetW3IDSecurityV1PublicKeyPem().Get()
		}
		p.Discoverable = asPerson.GetTootDiscoverable().Get()
		if v, ok := asPerson.GetUnknownProperties()["manuallyApprovesFollowers"]; ok {
			p.ManuallyApprovesFollowers = v.(bool)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	m, err := junk.FromBytes(b)
	if err != nil {
		return nil, err
	}
	if err := resolver.Resolve(context.Background(), m); err != nil {
		return nil, err
	}
	return p, nil
}

type OrderedCollection struct {
	ID           string              `json:"id,omitempty"`
	TotalItems   int                 `json:"totalItems,omitempty"`
	First        string              `json:"first,omitempty"`
	Last         string              `json:"last,omitempty"`
	OrderedItems []ActivityPubObject `json:"orderedItems,omitempty"`
}

func (o *OrderedCollection) GetID() string {
	return o.ID
}

func (o *OrderedCollection) GetType() string {
	return "OrderedCollection"
}

func (o *OrderedCollection) ToBytes() ([]byte, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	j, err := junk.FromBytes(b)
	if err != nil {
		return nil, err
	}
	j["type"] = o.GetType()
	deserialized, err := typeorderedcollection.DeserializeOrderedCollection(j, map[string]string{})
	if err != nil {
		return nil, err
	}
	m, err := streams.Serialize(deserialized)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type ActivityType string

const (
	FollowActivityType   ActivityType = "Follow"
	UndoActivityType     ActivityType = "Undo"
	AcceptActivityType   ActivityType = "Accept"
	RejectActivityType   ActivityType = "Reject"
	CreateActivityType   ActivityType = "Create"
	UpdateActivityType   ActivityType = "Update"
	DeleteActivityType   ActivityType = "Delete"
	AddActivityType      ActivityType = "Add"
	RemoveActivityType   ActivityType = "Remove"
	AnnounceActivityType ActivityType = "Announce"
)

type Activity struct {
	ID     string
	Type   ActivityType
	Actor  string
	Object ActivityPubObject
}

func (a *Activity) GetID() string {
	return a.ID
}

func (a *Activity) GetType() string {
	return string(a.Type)
}

func (a *Activity) ToBytes() ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	j, err := junk.FromBytes(b)
	if err != nil {
		return nil, err
	}
	j["type"] = a.GetType()
	deserialized, err := typeactivity.DeserializeActivity(j, map[string]string{})
	if err != nil {
		return nil, err
	}
	m, err := streams.Serialize(deserialized)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func ActivityFromBytes(b []byte) (*Activity, error) {
	j, err := junk.FromBytes(b)
	if err != nil {
		return nil, err
	}
	var a *Activity
	resolver, err := streams.NewJSONResolver(func(ctx context.Context, t vocab.ActivityStreamsActivity) error {
		a = &Activity{}
		a.ID = t.GetJSONLDId().Get().String()
		a.Type = ActivityType(t.GetJSONLDType().Begin().GetXMLSchemaString())
		a.Actor = t.GetActivityStreamsActor().Begin().GetIRI().String()
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := resolver.Resolve(context.Background(), j); err != nil {
		return nil, err
	}
	return a, nil
}
