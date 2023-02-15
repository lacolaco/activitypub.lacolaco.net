package ap

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"humungus.tedunangst.com/r/webs/httpsig"
)

type Person struct {
	*model.LocalUser
	baseURI string
}

func NewPerson(u *model.LocalUser, baseURI string) *Person {
	return &Person{LocalUser: u, baseURI: baseURI}
}

func (p *Person) GetID() string {
	return fmt.Sprintf("%s/users/%s", p.baseURI, p.LocalUser.UID)
}

func (p *Person) PubkeyID() string {
	return fmt.Sprintf("%s#%s", p.GetID(), publicKeyIDSuffix)
}

func (p *Person) InboxURI() string {
	return fmt.Sprintf("%s/inbox", p.GetID())
}

func (p *Person) OutboxURI() string {
	return fmt.Sprintf("%s/outbox", p.GetID())
}

func (p *Person) FollowersURI() string {
	return fmt.Sprintf("%s/followers", p.GetID())
}

func (p *Person) FollowingURI() string {
	return fmt.Sprintf("%s/following", p.GetID())
}

func (p *Person) GetProfileURL() string {
	return fmt.Sprintf("%s/@%s", p.baseURI, p.LocalUser.ID)
}

func (p *Person) ToMap(publicKey *rsa.PublicKey) map[string]interface{} {
	id := p.GetID()
	publicKeyPem, err := httpsig.EncodeKey(publicKey)
	if err != nil {
		panic(err)
	}

	return map[string]interface{}{
		"@context":                  contextURIs,
		"type":                      "Person",
		"id":                        id,
		"name":                      p.Name,
		"preferredUsername":         p.PrefName,
		"summary":                   p.Description,
		"inbox":                     p.InboxURI(),
		"outbox":                    p.OutboxURI(),
		"followers":                 p.FollowersURI(),
		"following":                 p.FollowingURI(),
		"url":                       p.GetProfileURL(),
		"published":                 p.CreatedAt.Format(time.RFC3339),
		"discoverable":              true,
		"manuallyApprovesFollowers": false,
		"icon": map[string]interface{}{
			"type":      "Image",
			"mediaType": p.Icon.MediaType,
			"url":       p.Icon.URL,
		},
		"publicKey": map[string]interface{}{
			"id":           p.PubkeyID(),
			"owner":        id,
			"publicKeyPem": publicKeyPem,
		},
	}
}
