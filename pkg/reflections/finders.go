package reflections

import (
	"fmt"
	"net/http"
)

func FindReflectedQueryValues(url string, res *http.Response, minLength uint) ([]*Reflection, error) {
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
