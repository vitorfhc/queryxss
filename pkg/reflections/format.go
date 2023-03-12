package reflections

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
