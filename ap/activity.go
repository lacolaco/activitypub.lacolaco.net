package ap

import (
	"encoding/json"
	"time"

	"humungus.tedunangst.com/r/webs/junk"
)

type ActivityType = ObjectType

type Activity struct {
	Context   interface{}  `json:"@context,omitempty"`
	ID        IRI          `json:"id,omitempty"`
	Type      ActivityType `json:"type,omitempty"`
	Actor     ObjectOrLink `json:"actor,omitempty"`
	Object    Object       `json:"object,omitempty"`
	Published time.Time    `json:"published,omitempty"`
	To        []Item       `json:"to,omitempty"`
}

var _ ActivityObject = (*Activity)(nil)

func (a *Activity) GetID() IRI {
	return a.ID
}

func (a *Activity) GetType() ObjectType {
	return a.Type
}

func (a *Activity) GetActor() ObjectOrLink {
	return a.Actor
}

func (a *Activity) GetObject() Object {
	return a.Object
}

func (a *Activity) MarshalJSON() ([]byte, error) {
	type temp Activity
	v := temp(*a)
	return json.Marshal(v)
}

func (a *Activity) UnmarshalJSON(data []byte) error {
	m, err := junk.FromBytes(data)
	if err != nil {
		return err
	}
	a.Context = m["@context"]
	if id, ok := m.GetString("id"); ok {
		a.ID = IRI(id)
	}
	if t, ok := m.GetString("type"); ok {
		a.Type = ActivityType(t)
	}
	if actor, ok := m.GetString("actor"); ok {
		a.Actor = IRI(actor)
	}
	if actor, ok := m.GetMap("actor"); ok {
		a.Actor = &Person{
			ID: IRI(actor["id"].(string)),
		}
	}
	if object, ok := m.GetString("object"); ok {
		a.Object = IRI(object)
	}
	if object, ok := m.GetMap("object"); ok {
		switch object["type"].(string) {
		case string(ObjectTypePerson):
			a.Object = &Person{
				ID: IRI(object["id"].(string)),
			}
		case string(ActivityTypeFollow):
			a.Object = &Activity{
				ID:   IRI(object["id"].(string)),
				Type: ActivityTypeFollow,
			}
		}
	}
	return nil
}

const (
	ActivityTypeAccept   ActivityType = "Accept"
	ActivityTypeAdd      ActivityType = "Add"
	ActivityTypeAnnounce ActivityType = "Announce"
	ActivityTypeCreate   ActivityType = "Create"
	ActivityTypeDelete   ActivityType = "Delete"
	ActivityTypeFollow   ActivityType = "Follow"
	ActivityTypeLike     ActivityType = "Like"
	ActivityTypeReject   ActivityType = "Reject"
	ActivityTypeRemove   ActivityType = "Remove"
	ActivityTypeUndo     ActivityType = "Undo"
	ActivityTypeUpdate   ActivityType = "Update"
)

func NewActivityAccept() *Activity {
	return &Activity{
		Type: ActivityTypeAccept,
	}
}

func NewActivityFollow() *Activity {
	return &Activity{
		Type: ActivityTypeFollow,
	}
}

func NewActivityUndo() *Activity {
	return &Activity{
		Type: ActivityTypeUndo,
	}
}
