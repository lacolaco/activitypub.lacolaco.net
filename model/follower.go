package model

import (
	"time"
)

type Follower struct {
	ID        string        `firestore:"-"`
	CreatedAt time.Time     `firestore:"created_at"`
	Status    AttemptStatus `firestore:"status"`
	UserID    string        `firestore:"user_id"`
}

func NewFollower(userID string, status AttemptStatus) *Follower {
	return &Follower{
		Status: status,
		UserID: userID,
	}
}

func (f *Follower) GetDocID() string {
	return f.ID
}
