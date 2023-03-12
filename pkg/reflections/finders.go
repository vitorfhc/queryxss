package reflections

import (
	"fmt"
	"strings"
)

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
