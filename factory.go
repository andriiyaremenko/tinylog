package tinylog

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/andriiyaremenko/tinylog/formatters"
)

// Returns new instance of `LoggerFactory` based on `out` and `formatter`
func NewLoggerFactory(out io.Writer, formatter formatters.LogFormatter) LoggerFactory {
	return &tinyLoggerFactory{
		out:       out,
		loggers:   make(map[context.Context]Logger),
		logLevel:  Info,
		formatter: formatter,
	}
}

// Returns new instance of `LoggerFactory` based on `os.Stderr` as `out` and `formatters.Default()` as `formatter`
func NewDefaultLoggerFactory() LoggerFactory {
	return NewLoggerFactory(os.Stderr, formatters.Default())
}

type tinyLoggerFactory struct {
	mu        sync.Mutex
	out       io.Writer
	loggers   map[context.Context]Logger
	logLevel  int
	formatter formatters.LogFormatter
}

func (tlf *tinyLoggerFactory) GetLogger(ctx context.Context) Logger {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	_, ok := tlf.loggers[ctx]

	if !ok {
		go func() {
			<-ctx.Done()
			tlf.mu.Lock()
			delete(tlf.loggers, ctx)
			tlf.mu.Unlock()
		}()
	}

	_, ok = tlf.loggers[ctx]

	if !ok {
		l := NewLogger(tlf.out, tlf.formatter)
		l.SetLogLevel(tlf.logLevel)
		tlf.loggers[ctx] = l
	}

	return tlf.loggers[ctx]
}

func (tlf *tinyLoggerFactory) SetLogLevel(level int) {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()
	tlf.logLevel = level
	for _, l := range tlf.loggers {
		l.SetLogLevel(level)
	}
}
