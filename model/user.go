package model

import "time"

type User struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        struct {
		URL       string `json:"url"`
		MediaType string `json:"media_type"`
	} `json:"icon"`
}

type Follower struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
