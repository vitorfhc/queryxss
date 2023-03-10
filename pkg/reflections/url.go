package reflections

import (
	"net/url"
	"strings"
)

// GetQueryValues returns a list of query values from a URL.
func GetQueryValues(u string) ([]string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	values := []string{}
	for _, v := range up.Query() {
		values = append(values, v...)
	}
	return values, nil
}

// GetQueryKeyValues returns a map of query keys and values from a URL.
func GetQueryKeyValues(u string) (map[string][]string, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return up.Query(), nil
}

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

// HasHtmlCharsInQueryValues returns true if the URL contains HTML special characters in the query values.
func HasHtmlCharsInQueryValues(u string) bool {
	values, err := GetQueryValues(u)
	if err != nil {
		return false
	}
	htmlChars := GetHtmlSpecialChars()
	for _, value := range values {
		found := FindAny(value, htmlChars, false)
		if found {
			return true
		}
	}
	return false
}

func RemoveHtmlCharsFromQueryValues(u string) (string, error) {
	query, err := GetQueryKeyValues(u)
	if err != nil {
		return "", err
	}
	if len(query) == 0 {
		return u, nil
	}

	htmlSpecialChars := GetHtmlSpecialChars()
	newQuery := map[string][]string{}
	for key, values := range query {
		newValues := []string{}
		for _, value := range values {
			newValue := ReplaceAll(value, htmlSpecialChars, RandomAlphaString(2))
			newValues = append(newValues, newValue)
		}
		newQuery[key] = newValues
	}

	return MergeQuery(u, newQuery)
}

func ReplaceAll(s string, old []string, new string) string {
	for _, o := range old {
		s = strings.ReplaceAll(s, o, new)
	}
	return s
}
