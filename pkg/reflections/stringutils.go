package reflections

import (
	"fmt"
	"io"
	"math/rand"
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

// ReflectionToString returns a string representation of a Reflection.
func ReflectionToString(r *Reflection, nocolor bool) string {
	toFormat := "[%s] %s reflected on %s [%s]"
	what := "%s (%s)"
	var msg string
	if nocolor {
		what = fmt.Sprintf(what, r.What, r.WhatName)
		msg = fmt.Sprintf(toFormat, r.Severity.String(), what, r.Where, r.Url)
	} else {
		severity := AddSeverityColor(r.Severity)
		what = fmt.Sprintf(what, r.What, r.WhatName)
		what = Bold(what)
		where := Bold(r.Where.String())
		msg = fmt.Sprintf("[%s] %s reflected on %s [%s]", severity, what, where, r.Url)
	}
	return msg
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

func ReplaceAll(s string, old []string, new string) string {
	for _, o := range old {
		s = strings.ReplaceAll(s, o, new)
	}
	return s
}
