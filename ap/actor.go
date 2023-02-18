package ap

type Actor interface {
	Object
	GetInbox() IRI
	GetOutbox() ObjectOrLink
	GetFollowers() ObjectOrLink
	GetFollowing() ObjectOrLink
}
