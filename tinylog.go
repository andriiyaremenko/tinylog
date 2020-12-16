package tinylog

import "context"

const (
	Trace int = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

type LogLevelSetter interface {
	SetLogLevel(level int)
}

type FixedLevelLogger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Logger interface {
	LogLevelSetter

	GetFixedLevel(level int) FixedLevelLogger
	AddTag(key string, value ...string)

	Printf(level int, format string, v ...interface{})
	Println(level int, v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
}

type LoggerFactory interface {
	LogLevelSetter
	GetLogger(ctx context.Context) Logger
}
