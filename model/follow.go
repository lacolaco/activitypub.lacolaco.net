package model

import (
	"encoding/json"
	"time"
)

type RemoteUser struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (u *RemoteUser) GetID() string {
	return u.ID
}

func (u *RemoteUser) ToMap() (map[string]interface{}, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	var h map[string]interface{}
	if err := json.Unmarshal(b, &h); err != nil {
		return nil, err
	}
	return h, nil
}
