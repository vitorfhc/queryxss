package reflections

// HtmlSpecialChars returns a list of HTML special characters.
// Since golang does not support constant arrays, we use a function instead.
func HtmlSpecialChars() []string {
	return []string{
		"\"",
		"'",
		"<",
		">",
	}
}
