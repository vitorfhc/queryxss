package reflections

import "net/http"

// ScanHttpClient is the interface that every HTTP client used by the scanners must implement.
// It's used so any HTTP client can be used by the scanners.
type ScanHttpClient interface {
	Get(url string) (*http.Response, error)
}
