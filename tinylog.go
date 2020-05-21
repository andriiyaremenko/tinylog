package tinylog

const (
	Debug = iota
	Info
	Warn
	Err
	Fatal
)

type LogLevelSetter interface {
	SetLogLevel(lovel int)
}

type TinyLogger interface {
	LogLevelSetter
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Err(v ...interface{})
	Errf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type TinyLoggerFactory interface {
	LogLevelSetter
	GetLogger(module string) TinyLogger
}
