package activitypub

type Object struct {
	Context any          `json:"@context,omitempty"`
	ID      string       `json:"id,omitempty"`
	Type    ActivityType `json:"type,omitempty"`
	Icon    Icon         `json:"icon,omitempty"`
}

type Icon struct {
	Type      string `json:"type,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	URL       string `json:"url,omitempty"`
}
