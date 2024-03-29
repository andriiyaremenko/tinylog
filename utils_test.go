package tinylog

import (
	"bytes"
	"context"
	"sync"

	"github.com/andriiyaremenko/tinylog/formatters"
)

func getLoggerFactory() (LoggerFactory, *bytes.Buffer) {
	b := new(bytes.Buffer)
	lf := NewLoggerFactory(DestinationFunc(b, formatters.Default(), Info))

	return lf, b
}

func getLogger() (Logger, []Destination, *bytes.Buffer) {
	ctx := context.TODO()
	lf, b := getLoggerFactory()
	l := lf.GetLogger(ctx, AllDestinations(lf)...)

	return l, AllDestinations(lf), b
}

type concurrentWriter struct {
	b  *bytes.Buffer
	mu sync.Mutex
}

func (cw *concurrentWriter) Write(p []byte) (n int, err error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.b.Write(p)
}

func (cw *concurrentWriter) String() string {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.b.String()
}
