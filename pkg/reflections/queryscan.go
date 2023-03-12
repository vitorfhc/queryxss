package reflections

import (
	"fmt"
	"net/url"
	"strings"
)

func QueryScan(httpClient ScanHttpClient, originalUrl string, minLength uint) ([]*Reflection, error) {
	originalUrlParsed, err := url.Parse(originalUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse original url: %v", err)
	}

	query := originalUrlParsed.Query()
	if len(query) == 0 {
		return []*Reflection{}, nil
	}

	for key, values := range query {
		for _, value := range values {
			if value == "" {
				token := RandomAlphaString(8)
				query.Set(key, token)
			}
		}
	}

	originalUrlParsed.RawQuery = query.Encode()
	response, err := httpClient.Get(originalUrlParsed.String())
	if err != nil {
		return nil, fmt.Errorf("failed to GET url: %v", err)
	}
	defer response.Body.Close()

	bodyAsString, err := ReaderToString(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	keyToReflections := map[string]*Reflection{}
	for key, values := range query {
		for _, value := range values {
			if len(value) < int(minLength) {
				continue
			}
			if _, ok := keyToReflections[key]; ok {
				continue
			}
			if Contains(bodyAsString, value, false) {
				keyToReflections[key] = &Reflection{
					Url:      originalUrlParsed.String(),
					Severity: "info",
					What:     fmt.Sprintf("query parameter [%s] with value [%s]", key, value),
					Where:    "body",
				}
			}
		}
	}

	originalQuery := CopyQuery(query)
	htmlSChars := GetHtmlSpecialChars()
	htmlSCharsStr := strings.Join(htmlSChars, "")
	tokenLen := 6
	for key := range keyToReflections {
		query = CopyQuery(originalQuery)
		tokenBegin := RandomAlphaString(tokenLen)
		tokenEnd := RandomAlphaString(tokenLen)
		newValue := tokenBegin + htmlSCharsStr + tokenEnd
		query.Set(key, newValue)
		originalUrlParsed.RawQuery = query.Encode()
		response, err := httpClient.Get(originalUrlParsed.String())
		if err != nil {
			return nil, fmt.Errorf("failed to GET url: %v", err)
		}
		defer response.Body.Close()
		bodyAsString, err := ReaderToString(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		founds := FindAllBetween(bodyAsString, tokenBegin, tokenEnd)
		for _, found := range founds {
			if FindAny(found, htmlSChars, false) {
				keyToReflections[key] = &Reflection{
					Url:      originalUrlParsed.String(),
					Severity: "medium",
					What:     fmt.Sprintf("query parameter [%s] with value [%s]", key, newValue),
					Where:    fmt.Sprintf("body [%s]", found),
				}
			}
		}
	}

	results := []*Reflection{}
	for _, reflection := range keyToReflections {
		results = append(results, reflection)
	}

	return results, nil
}
