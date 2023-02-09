package activitypub

const (
	ActivityPubContext = "https://www.w3.org/ns/activitystreams"
)

type ActivityType string

const (
	ActivityTypeAccept ActivityType = "Accept"
	ActivityTypeFollow ActivityType = "Follow"
	ActivityTypeUndo   ActivityType = "Undo"
	ActivityTypePerson ActivityType = "Person"
)

type Activity struct {
	Context any          `json:"@context,omitempty"`
	ID      string       `json:"id,omitempty"`
	Type    ActivityType `json:"type,omitempty"`
	Icon    Icon         `json:"icon,omitempty"`
	Object  Object       `json:"object,omitempty"`
	Actor   Actor        `json:"actor,omitempty"`
}

func (a *Activity) ToObject() Object {
	return Object{
		Context: a.Context,
		ID:      a.ID,
		Type:    a.Type,
		Icon:    a.Icon,
	}
}
