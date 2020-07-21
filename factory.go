package tinylog

import (
	"context"
	"io"
	"sync"
)

func NewTinyLoggerFactory(out io.Writer, format format, timeFormat string) LoggerFactory {
	return &tinyLoggerFactory{
		out:        out,
		loggers:    make(map[context.Context]map[string]Logger),
		logLevel:   Info,
		format:     format,
		timeFormat: timeFormat,
	}
}

type tinyLoggerFactory struct {
	mu         sync.Mutex
	out        io.Writer
	loggers    map[context.Context]map[string]Logger
	logLevel   logLevel
	format     format
	timeFormat string
}

func (tlf *tinyLoggerFactory) GetLogger(ctx context.Context, module string) Logger {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	_, ok := tlf.loggers[ctx]

	if !ok {
		tlf.loggers[ctx] = make(map[string]Logger)

		go func() {
			<-ctx.Done()
			tlf.mu.Lock()
			delete(tlf.loggers, ctx)
			tlf.mu.Unlock()
		}()
	}

	_, ok = tlf.loggers[ctx][module]

	if !ok {
		l := NewTinyLogger(tlf.out, tlf.format, module, tlf.timeFormat)
		l.SetLogLevel(tlf.logLevel)
		tlf.loggers[ctx][module] = l
	}

	return tlf.loggers[ctx][module]
}

func (tlf *tinyLoggerFactory) SetLogLevel(level logLevel) {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	tlf.logLevel = level
	for _, c := range tlf.loggers {
		for _, l := range c {
			l.SetLogLevel(level)
		}
	}
}
