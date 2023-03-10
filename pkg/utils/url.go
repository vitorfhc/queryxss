package utils

import "net/url"

// AddScheme adds the https scheme to a URL if it does not have one.
func AddScheme(u string) (string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	if up.Scheme == "" {
		up.Scheme = "https"
	}

	return up.String(), nil
}
