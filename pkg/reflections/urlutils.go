package reflections

import (
	"net/url"
)

// MergeQuery merges a map of query values into a URL.
// If the key already exists, it will be overwritten.
func MergeQuery(u string, q map[string][]string) (string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	qs := up.Query()
	for k, v := range q {
		qs[k] = v
	}

	up.RawQuery = qs.Encode()
	return up.String(), nil
}

// AddSchema adds the https scheme to a URL if it does not have one.
func AddSchema(u string) (string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	if up.Scheme == "" {
		up.Scheme = "https"
	}

	return up.String(), nil
}

// CopyQuery returns a copy of a url.Values.
func CopyQuery(src url.Values) url.Values {
	dst := make(url.Values)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
