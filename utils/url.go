package utils

import (
	"net/url"
	"path"
)

//JoinURL to join url
func JoinURL(base string, pths ...string) string {
	u, err := url.Parse(base)
	HandleError(err)
	for _, p := range pths {
		u.Path = path.Join(u.Path, p)
	}
	return u.String()
}
