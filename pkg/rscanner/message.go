package rscanner

import "fmt"

// SuccessMessage returns the success message of the scanner.
// It adds the severity color to the message if nocolor is false.
func SuccessMessage(info *ScannerInfo, url string, nocolor bool) string {
	severity := info.Severity
	id := info.ID
	if !nocolor {
		severity = AddSeverityColor(severity)
		id = Bold(id)
	}
	msg := fmt.Sprintf("[%s][%s] %s [%s]\n", id, severity, info.SuccessMessage, url)
	return msg
}

// Bold returns the string s in bold.
func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

// AddSeverityColor returns the string severity in the color
// corresponding to its severity.
// Info: blue
// Low: green
// Medium: yellow
// High: red
// Critical: magenta
func AddSeverityColor(severity string) string {
	switch severity {
	case "info":
		return "\033[1;34m" + severity + "\033[0m"
	case "low":
		return "\033[1;32m" + severity + "\033[0m"
	case "medium":
		return "\033[1;33m" + severity + "\033[0m"
	case "high":
		return "\033[1;31m" + severity + "\033[0m"
	case "critical":
		return "\033[1;35m" + severity + "\033[0m"
	default:
		return severity
	}
}
