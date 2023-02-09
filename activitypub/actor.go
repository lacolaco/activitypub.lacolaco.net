package activitypub

type Actor struct {
	Context           any          `json:"@context,omitempty"`
	ID                string       `json:"id,omitempty"`
	Type              ActivityType `json:"type,omitempty"`
	Name              string       `json:"name,omitempty"`
	Icon              Icon         `json:"icon,omitempty"`
	PreferredUsername string       `json:"preferredUsername,omitempty"`
	Summary           string       `json:"summary,omitempty"`
	Inbox             string       `json:"inbox,omitempty"`
	Outbox            string       `json:"outbox,omitempty"`
	URL               string       `json:"url,omitempty"`
	PublicKey         PublicKey    `json:"publicKey,omitempty"`
}
