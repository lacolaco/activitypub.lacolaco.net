package ap

var (
	systemActor = &actor{ID: "https://activitypub.lacolaco.net"}
)

type Actor interface {
	GetID() string
}

type actor struct {
	ID string
}

func (a *actor) GetID() string {
	return a.ID
}
