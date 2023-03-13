package reflections

type SeverityType int

const (
	SeverityInfo = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

const (
	SeverityInfoString     = "info"
	SeverityLowString      = "low"
	SeverityMediumString   = "medium"
	SeverityHighString     = "high"
	SeverityCriticalString = "critical"
)

func (s SeverityType) String() string {
	switch s {
	case SeverityInfo:
		return SeverityInfoString
	case SeverityLow:
		return SeverityLowString
	case SeverityMedium:
		return SeverityMediumString
	case SeverityHigh:
		return SeverityHighString
	case SeverityCritical:
		return SeverityCriticalString
	}
	return "UNKNOWN"
}

type Reflection struct {
	Url      string
	Severity SeverityType
	What     string
	Where    string
}

// ScanFunc is the function signature for the scan functions.
type ScanFunc func(client ScanHttpClient, url string, minLength uint) ([]*Reflection, error)
