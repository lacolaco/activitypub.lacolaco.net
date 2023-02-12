package ap

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"humungus.tedunangst.com/r/webs/httpsig"
)

type Person struct {
	*model.User
	baseURI string
	key     *rsa.PublicKey
}

func NewPerson(u *model.User, baseURI string, publicKey *rsa.PublicKey) *Person {
	return &Person{User: u, baseURI: baseURI, key: publicKey}
}

func (p *Person) ID() string {
	return fmt.Sprintf("%s/users/%s", p.baseURI, p.User.ID)
}

func (p *Person) PubkeyID() string {
	return fmt.Sprintf("%s#%s", p.ID(), DefaultPublicKeyID)
}

func (p *Person) InboxURI() string {
	return fmt.Sprintf("%s/inbox", p.ID())
}

func (p *Person) OutboxURI() string {
	return fmt.Sprintf("%s/outbox", p.ID())
}

func (p *Person) FollowersURI() string {
	return fmt.Sprintf("%s/followers", p.ID())
}

func (p *Person) FollowingURI() string {
	return fmt.Sprintf("%s/following", p.ID())
}

func (p *Person) AsMap() map[string]interface{} {
	id := p.ID()
	publicKeyPem, err := httpsig.EncodeKey(p.key)
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
		"url":                       id,
		"published":                 p.CreatedAt.Format(time.RFC3339),
		"updated":                   p.UpdatedAt.Format(time.RFC3339),
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
