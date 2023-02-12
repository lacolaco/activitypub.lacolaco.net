package webfinger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

const (
	WebFingerPath = "/.well-known/webfinger"
)

// "@username@domain" 形式のアカウント名からwebfingerを通じてPersonのURIを取得する
func ResolveAccountURI(ctx context.Context, resource string) (string, error) {
	re := regexp.MustCompile(`^@(.+)@(.+)$`)
	captured := re.FindStringSubmatch(resource)
	if len(captured) != 3 {
		return "", fmt.Errorf("invalid resource: %s", resource)
	}
	username := captured[1]
	domain := captured[2]

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s%s", domain, WebFingerPath), nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("resource", fmt.Sprintf("acct:%s@%s", username, domain))
	req.URL.RawQuery = q.Encode()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	defer resp.Body.Close()
	var o WebfingerObject
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return "", err
	}
	if len(o.Links) == 0 {
		return "", nil
	}
	for _, link := range o.Links {
		if link.Rel == "self" && link.Type == "application/activity+json" {
			return link.Href, nil
		}
	}
	return "", nil
}
