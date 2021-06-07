package tinylog

import (
	"context"
	"fmt"
	"sync"
)

// Returns new instance of LoggerFactory based on out and formatter.
func NewLoggerFactory(destinations ...Destination) LoggerFactory {
	if len(destinations) == 0 {
		panic("no destination was provided for LoggerFactory")
	}

	check := make(map[string]struct{})
	for _, destFunc := range destinations {
		dest := destFunc()
		if _, ok := check[dest.ID()]; ok {
			panic(fmt.Sprintf("destination %s was provided more than once", dest.ID()))
		}
	}

	return &tinyLoggerFactory{
		loggers:      make(map[context.Context]Logger),
		destinations: destinations}
}

// Returns new instance of LoggerFactory with DefaultDestination.
func DefaultLoggerFactory() LoggerFactory {
	return NewLoggerFactory(DefaultDestination)
}

type tinyLoggerFactory struct {
	mu sync.Mutex

	loggers      map[context.Context]Logger
	destinations []Destination
}

func (tlf *tinyLoggerFactory) Destinations() []Destination {
	return tlf.destinations
}

func (tlf *tinyLoggerFactory) GetLogger(ctx context.Context, destinations ...Destination) Logger {
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
		if len(destinations) == 0 {
			destinations = AllDestinations(tlf)
		}

		l := NewLogger(destinations...)
		tlf.loggers[ctx] = l
	}

	return tlf.loggers[ctx]
}

func (tlf *tinyLoggerFactory) SetLogLevel(level int, destinations ...Destination) {
	tlf.mu.Lock()
	defer tlf.mu.Unlock()

	for _, l := range tlf.loggers {
		l.SetLogLevel(level, destinations...)
	}
}
