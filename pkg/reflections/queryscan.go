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

	// TODO(vitorfhc): I think this could be improved.
	// Example: domain.tld/?q&q=2 would be replaced by domain.tld/?q=RANDOM
	// We are loosing the second q=2, which could be a reflection.
	for key, values := range query {
		for _, value := range values {
			if value == "" {
				token := RandomAlphaString(8)
				query.Set(key, token)
			}
		}
	}

	originalUrlParsed.RawQuery = query.Encode()
	modifiedUrl := originalUrlParsed.String()
	response, err := httpClient.Get(modifiedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to GET url: %v", err)
	}
	defer response.Body.Close()

	var responseHeader string
	for key, values := range response.Header {
		for _, value := range values {
			responseHeader += fmt.Sprintf("%s: %s\n", key, value)
		}
	}

	bodyAsString, err := ReaderToString(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	results := map[[2]string][]*Reflection{}
	for key, values := range query {
		id := [2]string{key, ""}
		results[id] = []*Reflection{}
		for _, value := range values {
			id = [2]string{key, value}
			results[id] = []*Reflection{}
		}
	}

	for key, values := range query {
		for _, value := range values {
			id := [2]string{key, value}
			if len(value) < int(minLength) {
				continue
			}
			if Contains(responseHeader, value, false) {
				results[id] = append(results[id], &Reflection{
					Url:      modifiedUrl,
					Severity: SeverityInfo,
					What:     WhatQueryValue,
					WhatName: value,
					Where:    WhereHeader,
				})
			}
			if Contains(bodyAsString, value, false) {
				results[id] = append(results[id], &Reflection{
					Url:      modifiedUrl,
					Severity: SeverityInfo,
					What:     WhatQueryValue,
					WhatName: value,
					Where:    WhereBody,
				})
			}
		}
	}

	for reflectionKey, reflections := range results {
		for _, reflection := range reflections {
			queryKey := reflectionKey[0]
			beginToken := RandomAlphaString(8)
			endToken := RandomAlphaString(8)
			// We don't want to search for header reflections
			if reflection.Where == WhereHeader {
				// We can add some CRLF here to make it more interesting
				continue
			} else {
				queryValue := strings.Join(GetHtmlSpecialChars(), "")
				newQuery := CopyQuery(query)
				newQuery.Set(queryKey, beginToken+queryValue+endToken)
				q := newQuery.Encode()
				originalUrlParsed.RawQuery = q
				newUrl := originalUrlParsed.String()
				response, err := httpClient.Get(newUrl)
				if err != nil {
					return nil, fmt.Errorf("failed to GET url: %v", err)
				}
				defer response.Body.Close()
				bodyAsString, err := ReaderToString(response.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %v", err)
				}
				findings := FindAllBetween(bodyAsString, beginToken, endToken)
				for _, finding := range findings {
					found := FindAny(finding, GetHtmlSpecialChars(), false)
					if found {
						reflection.Severity = SeverityMedium
						reflection.WhatName = queryValue
						reflection.Url = newUrl
					}
				}
			}
		}
	}

	resultsFlat := []*Reflection{}
	for _, reflections := range results {
		resultsFlat = append(resultsFlat, reflections...)
	}

	return resultsFlat, nil
}
