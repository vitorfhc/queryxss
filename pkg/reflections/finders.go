package reflections

import (
	"fmt"
	"net/http"
	"strings"
)

// FindReflectedValues requests a URL and returns a list of reflections.
// It searches for query values in the response body and headers.
func FindReflectedValues(url string, res *http.Response, minLength uint) ([]*Reflection, error) {
	values, err := GetQueryValues(url)
	if err != nil {
		return nil, err
	}

	bodyString, err := ReaderToString(res.Body)
	if err != nil {
		return nil, err
	}
	headers := res.Header
	results := []*Reflection{}

	for _, value := range values {
		if len(value) < int(minLength) {
			continue
		}
		for headerKey, headerValues := range headers {
			if Contains(headerKey, value, false) {
				result := &Reflection{
					Url:      url,
					Severity: "info",
					What:     value,
					Where:    fmt.Sprintf("header key %q", headerKey),
				}
				results = append(results, result)
			}
			for _, headerValue := range headerValues {
				if Contains(headerValue, value, false) {
					result := &Reflection{
						Url:      url,
						Severity: "info",
						What:     fmt.Sprintf("query value %q", value),
						Where:    fmt.Sprintf("header (%q) value", headerKey),
					}
					results = append(results, result)
				}
			}
		}

		if Contains(bodyString, value, false) {
			result := &Reflection{
				Url:      url,
				Severity: "info",
				What:     fmt.Sprintf("query value %q", value),
				Where:    "body",
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// FindBetween returns the substring between two substrings.
func FindBetween(s string, begin string, end string) (string, error) {
	beginIndex := strings.Index(s, begin)
	if beginIndex == -1 {
		return "", fmt.Errorf("begin string %q not found", begin)
	}
	beginIndex += len(begin)
	endIndex := strings.Index(s[beginIndex:], end)
	if endIndex == -1 {
		return "", fmt.Errorf("end string %q not found", end)
	}
	endIndex += len(s[:beginIndex])
	return s[beginIndex:endIndex], nil
}

// FindAllBetween returns all substrings between two substrings.
func FindAllBetween(s string, begin string, end string) []string {
	results := []string{}
	for {
		found, err := FindBetween(s, begin, end)
		if err != nil {
			break
		}
		results = append(results, found)
		s = s[strings.Index(s, end)+len(end):]
	}
	return results
}
