package reflections

type Reflection struct {
	Url      string
	Severity string
	What     string
	Where    string
}

// ScanFunc is the function signature for the scan functions.
type ScanFunc func(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error)
