package reflections

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
)

// ReaderToString reads all the data from an io.Reader and returns it as a string.
func ReaderToString(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Contains returns true if s contains any of the strings in subs.
// If caseSensitive is true, the comparison is case sensitive.
func Contains(s string, sub string, caseSensitive bool) bool {
	if !caseSensitive {
		s = strings.ToLower(s)
		sub = strings.ToLower(sub)
	}
	return strings.Contains(s, sub)
}

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

// ReflectionToString returns a string representation of a Reflection.
func ReflectionToString(r *Reflection, nocolor bool) string {
	var msg string
	if nocolor {
		msg = fmt.Sprintf("[%s] %s reflected on %s [%s]", r.Severity, r.What, r.Where, r.Url)
	} else {
		severity := AddSeverityColor(r.Severity)
		what := Bold(r.What)
		where := Bold(r.Where)
		msg = fmt.Sprintf("[%s] %s reflected on %s [%s]", severity, what, where, r.Url)
	}
	return msg
}

// Bold returns the string s in bold.
func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

// AddSeverityColor returns the string severity in the color
// corresponding to its severity.
func AddSeverityColor(severity string) string {
	switch severity {
	case "info": // blue
		return "\033[1;34m" + severity + "\033[0m"
	case "low": // green
		return "\033[1;32m" + severity + "\033[0m"
	case "medium": // yellow
		return "\033[1;33m" + severity + "\033[0m"
	case "high": // red
		return "\033[1;31m" + severity + "\033[0m"
	case "critical": // magenta
		return "\033[1;35m" + severity + "\033[0m"
	default:
		return severity
	}
}

// RandomAlphaString returns a random string of length n.
// The string is composed of lowercase and uppercase letters.
func RandomAlphaString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// FindAll returns a list of indices where substr is found in s.
func FindAll(s string, substr string, caseSensitive bool) []int {
	if !caseSensitive {
		s = strings.ToLower(s)
		substr = strings.ToLower(substr)
	}

	var indices []int
	for i := 0; ; i++ {
		i = strings.Index(s[i:], substr)
		if i == -1 {
			break
		}
		indices = append(indices, i)
		i = i + len(substr)
	}
	return indices
}

// FindAny returns true if s contains any of the strings in subs.
func FindAny(s string, subs []string, caseSensitive bool) bool {
	if !caseSensitive {
		s = strings.ToLower(s)
		for i, sub := range subs {
			subs[i] = strings.ToLower(sub)
		}
	}

	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
