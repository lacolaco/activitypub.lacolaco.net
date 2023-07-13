package model

import (
	"time"
)

type Follower struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"createdAt"`
}

func NewFollower(userID string, status AttemptStatus) *Follower {
	return &Follower{
		ID: userID,
	}
}
