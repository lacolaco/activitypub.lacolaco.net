package model

import (
	"time"
)

type UID string

type UserIcon struct {
	URL       string `json:"url" firestore:"url"`
	MediaType string `json:"media_type" firestore:"media_type"`
}

type LocalUser struct {
	UID         UID       `json:"uid" firestore:"-"`
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	PrefName    string    `json:"preferred_username" firestore:"preferred_username"`
	Description string    `json:"description" firestore:"description"`
	Icon        *UserIcon `json:"icon" firestore:"icon"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}

func (u *LocalUser) GetDocID() string {
	return string(u.UID)
}

type RemoteUser struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (u *RemoteUser) GetID() string {
	return u.ID
}
