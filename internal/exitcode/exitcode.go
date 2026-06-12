package exitcode

const (
	OK              = 0
	PolicyFailed    = 1
	ConfigError     = 2
	Environment     = 3
	ScannerFailed   = 4
	SafetyViolation = 5
	Interrupted     = 130
)
