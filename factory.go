package tinylog

import "io"

func NewTinyLoggerFactory(out io.Writer) TinyLoggerFactory {
	return &tinyLoggerFactory{
		out:     out,
		loggers: make(map[string]TinyLogger),
	}
}

type tinyLoggerFactory struct {
	out     io.Writer
	loggers map[string]TinyLogger
}

func (tlf *tinyLoggerFactory) GetLogger(module string) TinyLogger {
	l, ok := tlf.loggers[module]
	if ok {
		return l
	}
	l = NewTinyLogger(tlf.out, module)
	tlf.loggers[module] = l
	return l
}
