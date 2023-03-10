package reflections

// GetHtmlSpecialChars returns a list of HTML special characters.
// Since golang does not support constant arrays, we use a function instead.
func GetHtmlSpecialChars() []string {
	return []string{
		"\"",
		"'",
		"<",
		">",
	}
}
