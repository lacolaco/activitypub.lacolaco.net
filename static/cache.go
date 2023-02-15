package static

import "regexp"

func DetectCacheControl(filename string) string {
	switch {
	case regexp.MustCompile(`.+\.[0-9a-f]{16,}\.(css|js)$`).Match([]byte(filename)):
		return "public, max-age=31536000, immutable"
	}
	return "no-cache"
}
