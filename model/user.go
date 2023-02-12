package model

import (
	"time"
)

type UserIcon struct {
	URL       string `firestore:"url"`
	MediaType string `firestore:"media_type"`
}

type LocalUser struct {
	ID          string    `firestore:"id"`
	Name        string    `firestore:"name"`
	PrefName    string    `firestore:"preferred_username"`
	Description string    `firestore:"description"`
	Icon        *UserIcon `firestore:"icon"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

type RemoteUser struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (u *RemoteUser) GetID() string {
	return u.ID
}
