package model

import (
	"encoding/json"
	"time"
)

type RemoteUser struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
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

func NewRemoteUserFromMap(v map[string]interface{}) (*RemoteUser, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var u *RemoteUser
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, err
	}
	return u, nil
}
