package reflections

import "net/http"

// ScanHttpClient is the interface used by the scan functions to request URLs.
type ScanHttpClient interface {
	Get(url string) (*http.Response, error)
}
