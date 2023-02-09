package activitypub

type PublicKey struct {
	Context      any    `json:"@context,omitempty"`
	Type         string `json:"type,omitempty"`
	ID           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}
