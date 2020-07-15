package tinylog

import (
	"io"
	"sync"
)

func NewTinyLoggerFactory(out io.Writer, format format, timeFormat string) LoggerFactory {
	return &tinyLoggerFactory{
		out:        out,
		loggers:    make(map[string]Logger),
		logLevel:   Info,
		format:     format,
		timeFormat: timeFormat,
	}
}

type tinyLoggerFactory struct {
	mu         sync.Mutex
	out        io.Writer
	loggers    map[string]Logger
	logLevel   logLevel
	format     format
	timeFormat string
}

func (tlf *tinyLoggerFactory) GetLogger(module string) Logger {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	l, ok := tlf.loggers[module]
	if ok {
		return l
	}
	l = NewTinyLogger(tlf.out, tlf.format, module, tlf.timeFormat)
	l.SetLogLevel(tlf.logLevel)
	tlf.loggers[module] = l
	return l
}

func (tlf *tinyLoggerFactory) SetLogLevel(level logLevel) {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	tlf.logLevel = level
	for _, l := range tlf.loggers {
		l.SetLogLevel(level)
	}
}
