package ap

import (
	"encoding/json"
	"time"

	"humungus.tedunangst.com/r/webs/junk"
)

const (
	ObjectTypePerson ObjectType = "Person"
)

type PublicKey struct {
	ID           IRI    `json:"id,omitempty"`
	Owner        IRI    `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Person struct {
	Context           interface{} `json:"@context,omitempty"`
	ID                IRI         `json:"id,omitempty"`
	Type              ObjectType  `json:"type,omitempty"`
	Name              string      `json:"name,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	Inbox             IRI         `json:"inbox,omitempty"`
	Outbox            IRI         `json:"outbox,omitempty"`
	Followers         IRI         `json:"followers,omitempty"`
	Following         IRI         `json:"following,omitempty"`
	Liked             IRI         `json:"liked,omitempty"`
	URL               string      `json:"url,omitempty"`
	Published         time.Time   `json:"published,omitempty"`
	Icon              *Image      `json:"icon,omitempty"`
	PublicKey         *PublicKey
	Attachment        []ActivityStreamsObject `json:"attachment,omitempty"`

	Discoverable              bool `json:"discoverable"`
	ManuallyApprovesFollowers bool `json:"manuallyApprovesFollowers"`
}

var _ ObjectOrLink = (*Person)(nil)
var _ Actor = (*Person)(nil)

func (p *Person) GetID() IRI {
	return p.ID
}

func (p *Person) GetType() ObjectType {
	return ObjectTypePerson
}

func (p *Person) GetLink() IRI {
	return p.GetID()
}

func (p *Person) IsLink() bool {
	return false
}

func (p *Person) IsObject() bool {
	return true
}

func (p *Person) GetInbox() IRI {
	return p.Inbox
}

func (p *Person) GetOutbox() ObjectOrLink {
	return p.Outbox
}

func (p *Person) GetFollowers() ObjectOrLink {
	return p.Followers
}

func (p *Person) GetFollowing() ObjectOrLink {
	return p.Following
}

func (p *Person) MarshalJSON() ([]byte, error) {
	type temp Person
	v := temp(*p)
	v.Type = p.GetType()
	return json.Marshal(v)
}

func (p *Person) UnmarshalJSON(data []byte) error {
	m, err := junk.FromBytes(data)
	if err != nil {
		return err
	}

	if v, ok := m["@context"]; ok {
		p.Context = v
	}
	if v, ok := m.GetString("id"); ok {
		p.ID = IRI(v)
	}
	if v, ok := m.GetString("name"); ok {
		p.Name = v
	}
	if v, ok := m.GetString("preferredUsername"); ok {
		p.PreferredUsername = v
	}
	if v, ok := m.GetString("summary"); ok {
		p.Summary = v
	}
	if v, ok := m.GetString("inbox"); ok {
		p.Inbox = IRI(v)
	}
	if v, ok := m.GetString("outbox"); ok {
		p.Outbox = IRI(v)
	}
	if v, ok := m.GetString("followers"); ok {
		p.Followers = IRI(v)
	}
	if v, ok := m.GetString("following"); ok {
		p.Following = IRI(v)
	}
	if v, ok := m.GetString("liked"); ok {
		p.Liked = IRI(v)
	}
	if v, ok := m.GetString("url"); ok {
		p.URL = v
	}
	if v, ok := m.GetString("published"); ok {
		published, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		p.Published = published
	}
	if v, ok := m.GetString("icon"); ok {
		p.Icon = &Image{
			URL: v,
		}
	}
	if v, ok := m.GetMap("icon"); ok {
		p.Icon = &Image{
			URL: v["url"].(string),
		}
		if mediaType, ok := v.GetString("mediaType"); ok {
			p.Icon.MediaType = mediaType
		}
	}
	if v, ok := m.GetMap("publicKey"); ok {
		p.PublicKey = &PublicKey{
			ID:           IRI(v["id"].(string)),
			Owner:        IRI(v["owner"].(string)),
			PublicKeyPem: v["publicKeyPem"].(string),
		}
	}
	if v, ok := m.GetArray("attachment"); ok {
		p.Attachment = make([]ActivityStreamsObject, 0, len(v))
		for _, a := range v {
			m := a.(junk.Junk)
			switch m["type"] {
			case string(ObjectTypePropertyValue):
				p.Attachment = append(p.Attachment, &PropertyValue{
					Name:  m["name"].(string),
					Value: m["value"].(string),
				})
			}
		}
	}
	if v, ok := m["discoverable"]; ok {
		p.Discoverable = v.(bool)
	}
	if v, ok := m["manuallyApprovesFollowers"]; ok {
		p.ManuallyApprovesFollowers = v.(bool)
	}

	return nil
}
