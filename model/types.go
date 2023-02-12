package model

type AttemptStatus = string

const (
	AttemptStatusCompleted AttemptStatus = "completed"
	AttemptStatusPending   AttemptStatus = "pending"
)
