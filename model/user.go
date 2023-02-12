package model

import (
	"time"
)

type UserIcon struct {
	URL       string `firestore:"url"`
	MediaType string `firestore:"media_type"`
}

type User struct {
	ID          string    `firestore:"id"`
	Name        string    `firestore:"name"`
	PrefName    string    `firestore:"preferred_username"`
	Description string    `firestore:"description"`
	Icon        UserIcon  `firestore:"icon"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}
