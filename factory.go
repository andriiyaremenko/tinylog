package tinylog

import (
	"io"
	"sync"
)

func NewTinyLoggerFactory(out io.Writer) TinyLoggerFactory {
	return &tinyLoggerFactory{
		out:      out,
		loggers:  make(map[string]TinyLogger),
		logLevel: Info,
	}
}

type tinyLoggerFactory struct {
	mu       sync.Mutex
	out      io.Writer
	loggers  map[string]TinyLogger
	logLevel LogLevel
}

func (tlf *tinyLoggerFactory) GetLogger(module string) TinyLogger {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	l, ok := tlf.loggers[module]
	if ok {
		return l
	}
	l = NewTinyLogger(tlf.out, module)
	l.SetLogLevel(tlf.logLevel)
	tlf.loggers[module] = l
	return l
}

func (tlf *tinyLoggerFactory) SetLogLevel(level LogLevel) {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	tlf.logLevel = level
	for _, l := range tlf.loggers {
		l.SetLogLevel(level)
	}
}
