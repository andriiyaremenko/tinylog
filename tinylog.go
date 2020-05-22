package tinylog

const (
	Debug = iota
	Info
	Warn
	Error
	Fatal
	None
)

type LogLevelSetter interface {
	SetLogLevel(level int)
}

type TinyLogger interface {
	LogLevelSetter
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

type TinyLoggerFactory interface {
	LogLevelSetter
	GetLogger(module string) TinyLogger
}
