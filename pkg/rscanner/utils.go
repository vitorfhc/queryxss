package rscanner

import (
	"net/url"
	"strings"
)

// CaseInsensitiveContains returns true if s contains any of the strings in subs.
// The comparison is case insensitive.
func CaseInsensitiveContains(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(s), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// GetURLQuery returns the query parameters of a URL.
// It returns an error if the URL is invalid.
func GetURLQuery(u string) (url.Values, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return up.Query(), nil
}
