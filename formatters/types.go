package formatters

type LogFormatter interface {
	GetOutput(level int, message string, tags map[string][]string, calldepth int) []byte
}
