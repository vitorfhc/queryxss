package reflections

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Reflection struct {
	Url      string
	Severity string
	What     string
	Where    string
}

// ScanFunc is the function signature for the scan functions.
type ScanFunc func(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error)

// SimpleScan requests the url with no modifications and searches for query values in the response body and headers.
// It returns all findings as Reflections.
func SimpleScan(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error) {
	logrus.Debug("running simple scan")
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reflections, err := FindReflectedValues(url, res, minLength)
	if err != nil {
		return nil, err
	}

	if len(reflections) > 0 || !HasHtmlCharsInQueryValues(url) {
		return reflections, nil
	}

	// We got no reflections and the url has HTML special characters in the query values.
	// Maybe we didn't get reflections because these characters were encoded.
	// Let's remove the HTML special characters and try again.
	logrus.Debug("running simple scan without HTML special characters")
	urlMod, err := RemoveHtmlCharsFromQueryValues(url)
	if err != nil {
		return nil, err
	}
	res, err = client.Get(urlMod)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reflections, err = FindReflectedValues(urlMod, res, minLength)
	if err != nil {
		return nil, err
	}

	return reflections, nil
}

// ReplaceValuesHtmlCharsScan requests the url with modified query values.
// It replaces the values with HTML special characters and searches for the modified values in the response body.
// If the modified values are found and were not encoded, it returns it as a Reflection.
func ReplaceValuesHtmlCharsScan(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error) {
	logrus.Debug("running postfix special chars values scan")

	query, err := GetQueryKeyValues(url)
	if err != nil {
		return nil, err
	}
	if len(query) == 0 {
		return []*Reflection{}, nil
	}

	htmlSpecialChars := []string{"<", ">", "\"", "'"}
	begin := RandomAlphaString(6)
	end := RandomAlphaString(6)
	htmlChars := strings.Join(htmlSpecialChars, "")
	payload := begin + "%s" + htmlChars + end

	newQuery := map[string][]string{}
	for key, values := range query {
		newValues := []string{}
		for _, value := range values {
			newValue := fmt.Sprintf(payload, value)
			newValues = append(newValues, newValue)
		}
		newQuery[key] = newValues
	}

	urlMod, err := MergeQuery(url, newQuery)
	if err != nil {
		return nil, err
	}

	res, err := client.Get(urlMod)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyString, err := ReaderToString(res.Body)
	if err != nil {
		return nil, err
	}

	reflections := []*Reflection{}
	results := FindAllBetween(bodyString, begin, end)
	reflectedChars := []string{}
	for _, result := range results {
		for _, htmlChar := range htmlSpecialChars {
			if strings.Contains(result, htmlChar) {
				reflectedChars = append(reflectedChars, htmlChar)
			}
		}
	}

	if len(reflectedChars) > 0 {
		joined := strings.Join(reflectedChars, ", ")
		reflection := &Reflection{
			Url:      urlMod,
			Severity: "medium",
			What:     fmt.Sprintf("special characters [%s]", joined),
			Where:    "body",
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}
