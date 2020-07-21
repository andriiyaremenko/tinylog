package tinylog

import "context"

type logLevel int

const (
	Debug logLevel = iota
	Info
	Warn
	Error
	Fatal
	None
)

type Stringer interface {
	String() string
}

type LogLevelSetter interface {
	SetLogLevel(level logLevel)
}

type Logger interface {
	LogLevelSetter

	AddTag(ctx context.Context, key string, value ...string)

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type LoggerFactory interface {
	LogLevelSetter
	GetLogger(module string) Logger
}
