package model

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"humungus.tedunangst.com/r/webs/httpsig"
)

const (
	publicKeyIDSuffix = "key"
)

type UID string

type UserIcon struct {
	URL       string `json:"url" firestore:"url"`
	MediaType string `json:"media_type" firestore:"media_type"`
}

type LocalUser struct {
	UID         UID       `json:"uid" firestore:"-"`
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	PrefName    string    `json:"preferred_username" firestore:"preferred_username"`
	Description string    `json:"description" firestore:"description"`
	Icon        *UserIcon `json:"icon" firestore:"icon"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}

func (u *LocalUser) GetDocID() string {
	return string(u.UID)
}

func (u *LocalUser) ToPerson(baseURI string, publicKey *rsa.PublicKey) *ap.Person {
	id := fmt.Sprintf("%s/users/%s", baseURI, u.UID)
	publicKeyPem, err := httpsig.EncodeKey(publicKey)
	if err != nil {
		panic(err)
	}

	p := &ap.Person{
		ID:                ap.IRI(id),
		Name:              u.Name,
		PreferredUsername: u.PrefName,
		Summary:           u.Description,
		Inbox:             ap.IRI(fmt.Sprintf("%s/inbox", id)),
		Outbox:            ap.IRI(fmt.Sprintf("%s/outbox", id)),
		Followers:         ap.IRI(fmt.Sprintf("%s/followers", id)),
		Following:         ap.IRI(fmt.Sprintf("%s/following", id)),
		Liked:             ap.IRI(fmt.Sprintf("%s/liked", id)),
		URL:               fmt.Sprintf("%s/@%s", baseURI, u.ID),
		Published:         u.CreatedAt,
		Icon: &ap.Image{
			URL:       u.Icon.URL,
			MediaType: u.Icon.MediaType,
		},
		PublicKey: &ap.PublicKey{
			ID:           ap.IRI(fmt.Sprintf("%s#%s", id, publicKeyIDSuffix)),
			Owner:        ap.IRI(id),
			PublicKeyPem: publicKeyPem,
		},
		Attachment: []ap.ActivityStreamsObject{
			&ap.PropertyValue{
				Name:  "Twitter",
				Value: "\u003ca href=\"https://twitter.com/laco2net\" rel=\"me nofollow noopener\" target=\"_blank\"\u003e\u003cspan class=\"invisible\"\u003ehttps://\u003c/span\u003e\u003cspan class=\"\"\u003etwitter.com/laco2net\u003c/span\u003e\u003cspan class=\"invisible\"\u003e\u003c/span\u003e\u003c/a\u003e",
			},
		},
		Discoverable:              true,
		ManuallyApprovesFollowers: false,
	}
	return p
}

type RemoteUser struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (u *RemoteUser) GetID() string {
	return u.ID
}
