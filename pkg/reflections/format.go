package reflections

// Bold returns the string s in bold.
func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

// AddSeverityColor returns the string severity in the color
// corresponding to its severity.
func AddSeverityColor(severity SeverityType) string {
	switch severity {
	case SeverityInfo: // blue
		return "\033[1;34m" + SeverityInfoString + "\033[0m"
	case SeverityLow: // green
		return "\033[1;32m" + SeverityLowString + "\033[0m"
	case SeverityMedium: // yellow
		return "\033[1;33m" + SeverityMediumString + "\033[0m"
	case SeverityHigh: // red
		return "\033[1;31m" + SeverityHighString + "\033[0m"
	case SeverityCritical: // magenta
		return "\033[1;35m" + SeverityCriticalString + "\033[0m"
	}
	return "UNKNOWN"
}
