package model

import (
	"encoding/json"
	"fmt"

	"github.com/lacolaco/activitypub.lacolaco.net/sign"
)

type UserIcon struct {
	URL       string `json:"url"`
	MediaType string `json:"media_type"`
}

type User struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	PrefName    string   `json:"preferred_username"`
	Description string   `json:"description"`
	Icon        UserIcon `json:"icon"`
}

func NewUserFromMap(v interface{}) (*User, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var user *User
	if err := json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetActivityPubID(baseURI string) string {
	return fmt.Sprintf("%s/users/%s", baseURI, u.ID)
}

func (u *User) GetPubkeyID(baseURI string) string {
	return fmt.Sprintf("%s#%s", u.GetActivityPubID(baseURI), sign.DefaultPublicKeyID)
}
