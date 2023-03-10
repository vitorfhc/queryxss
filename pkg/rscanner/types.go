package rscanner

import (
	"net/http"
)

// Input is the input for a scanner.
type Input struct {
	Url     string
	Headers map[string][]string
}

// ScannerInfo contains information about a scanner.
// Every scanner must have a ScannerInfo, because it is used to
// identify the scanner, display information about it, and build
// the final success message.
type ScannerInfo struct {
	// ID is the unique identifier of the scanner.
	// It must contain only lowercase letters and dashes.
	ID string
	// Description is a short description of the scanner.
	// It is used when listing available scanners.
	Description string
	// Severity is the severity of the vulnerability.
	// It must be one of the following: info, low, medium, high, critical.
	Severity string
	// SuccessMessage is the message to be displayed when the scanner
	// detects a vulnerability.
	// Example: "query reflection on headers detected"
	SuccessMessage string
}

// ReflectionScanner is the interface that every scanner must implement.
type ReflectionScanner interface {
	// Scan scans the input and returns true if the scanner
	// detects a vulnerability. It returns false otherwise.
	// If an error occurs, it returns false and the error.
	Scan(HttpRequester, *Input) (*Output, error)
	// GetInfo retrieves the ScannerInfo of the scanner.
	GetInfo() *ScannerInfo
}

// HttpRequester is the interface that every HTTP client used by the scanners must implement.
type HttpRequester interface {
	Get(url string) (*http.Response, error)
}

type Output struct {
	Success bool
	Info    *ScannerInfo
	Url     string
}
