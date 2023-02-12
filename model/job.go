package model

import "time"

type JobType = string

const (
	JobTypeFollowUser JobType = "follow_user"
)

type Job struct {
	ID        string      `firestore:"id"`
	Type      JobType     `firestore:"type"`
	CreatedAt time.Time   `firestore:"created_at"`
	CreatedBy string      `firestore:"created_by"`
	Target    interface{} `firestore:"target"`
}

func NewJob(id string, jobType JobType, userID string, target interface{}) *Job {
	return &Job{
		ID:        id,
		Type:      jobType,
		CreatedAt: time.Now(),
		CreatedBy: userID,
		Target:    target,
	}
}

func (j *Job) GetID() string {
	return j.ID
}
