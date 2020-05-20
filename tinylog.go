package tinylog

import ()

type TinyLogger interface {
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
	GetLogger(module string) TinyLogger
}
