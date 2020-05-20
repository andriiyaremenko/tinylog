package tinylog

import "io"

func NewTinyLoggerFabric(out io.Writer) TinyLoggerFabric {
	return &tinyLoggerFabric{
		out:     out,
		loggers: make(map[string]TinyLogger),
	}
}

type tinyLoggerFabric struct {
	out     io.Writer
	loggers map[string]TinyLogger
}

func (tlf *tinyLoggerFabric) GetLogger(module string) TinyLogger {
	l, ok := tlf.loggers[module]
	if ok {
		return l
	}
	l = NewTinyLogger(tlf.out, module)
	tlf.loggers[module] = l
	return l
}
