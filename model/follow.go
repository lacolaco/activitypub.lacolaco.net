package model

import (
	"encoding/json"
	"time"
)

type Follower struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Follower) ToMap() (map[string]interface{}, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	var h map[string]interface{}
	if err := json.Unmarshal(b, &h); err != nil {
		return nil, err
	}
	return h, nil
}

func NewFollowerFromMap(v map[string]interface{}) (*Follower, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var follower *Follower
	if err := json.Unmarshal(b, &follower); err != nil {
		return nil, err
	}
	return follower, nil
}
