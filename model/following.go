package model

import (
	"time"
)

type Following struct {
	ID        string        `firestore:"-"`
	CreatedAt time.Time     `firestore:"created_at"`
	Status    AttemptStatus `firestore:"status"`
	UserID    string        `firestore:"user_id"`
}

func NewFollowing(userID string, status AttemptStatus) *Following {
	return &Following{
		Status: status,
		UserID: userID,
	}
}

func (f *Following) GetDocID() string {
	return f.ID
}
