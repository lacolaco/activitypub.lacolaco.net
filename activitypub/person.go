package activitypub

type Person struct {
	Context           string `json:"@context,omitempty"`
	ID                string `json:"id,omitempty"`
	Type              string `json:"type,omitempty"`
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferredUsername,omitempty"`
	Summary           string `json:"summary,omitempty"`
	Inbox             string `json:"inbox,omitempty"`
	Outbox            string `json:"outbox,omitempty"`
	URL               string `json:"url,omitempty"`
	Icon              Icon   `json:"icon,omitempty"`
}

type Icon struct {
	Type      string `json:"type,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	URL       string `json:"url,omitempty"`
}
