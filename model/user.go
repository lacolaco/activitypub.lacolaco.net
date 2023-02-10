package model

type User struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        struct {
		URL       string `json:"url"`
		MediaType string `json:"media_type"`
	} `json:"icon"`
}
