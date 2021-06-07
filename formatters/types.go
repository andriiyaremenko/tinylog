package formatters

// Carries log message formatting and marshalling logic.
type LogFormatter interface {
	// Returns formatted log message in []byte.
	GetOutput(level int, message string, tags map[string][]string, calldepth int) []byte
}
