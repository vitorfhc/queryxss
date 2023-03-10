package reflections

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

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
