package model

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"

	goap "github.com/go-ap/activitypub"
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

func (u *User) ToPerson(baseUri string, publicKey *rsa.PublicKey) *goap.Person {
	apID := u.GetActivityPubID(baseUri)

	return &goap.Person{
		Type:              goap.PersonType,
		ID:                goap.IRI(apID),
		Name:              goap.DefaultNaturalLanguageValue(u.Name),
		PreferredUsername: goap.DefaultNaturalLanguageValue(u.PrefName),
		Summary:           goap.DefaultNaturalLanguageValue(u.Description),
		Inbox:             goap.IRI(fmt.Sprintf("%s/inbox", apID)),
		Outbox:            goap.IRI(fmt.Sprintf("%s/outbox", apID)),
		Followers:         goap.IRI(fmt.Sprintf("%s/followers", apID)),
		Following:         goap.IRI(fmt.Sprintf("%s/following", apID)),
		URL:               goap.IRI(fmt.Sprintf("%s/@%s", baseUri, u.ID)),
		Icon: &goap.Object{
			Type:      goap.ImageType,
			MediaType: goap.MimeType(u.Icon.MediaType),
			URL:       goap.IRI(u.Icon.URL),
		},
		PublicKey: goap.PublicKey{
			ID:           goap.ID(fmt.Sprintf("%s#%s", apID, sign.DefaultPublicKeyID)),
			Owner:        goap.IRI(apID),
			PublicKeyPem: sign.ExportPublicKey(publicKey),
		},
	}
}
