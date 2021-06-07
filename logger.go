package tinylog

import (
	"fmt"
	"os"
	"sync"

	"github.com/andriiyaremenko/tinylog/formatters"
)

// Returns new instance of Logger based on out and formatter.
func NewLogger(destinations ...Destination) Logger {
	if len(destinations) == 0 {
		panic("no destination was provided for Logger")
	}

	validated := make([]*destination, 0, 1)

	check := make(map[string]struct{})
	for _, destFunc := range destinations {
		dest := destFunc()
		if _, ok := check[dest.ID()]; ok {
			panic(fmt.Sprintf("destination %s was provided more than once", dest.ID()))
		}

		validated = append(validated, dest)
	}

	return &tinyLogger{
		tags:         make(map[string][]string),
		destinations: validated}
}

// Returns new instance of Logger with DefaultDestination.
func DefaultLogger() Logger {
	return NewLogger(DefaultDestination)
}

type fixedLevelLogger struct {
	l     *tinyLogger
	level int
}

func (fll *fixedLevelLogger) Printf(format string, v ...interface{}) {
	fll.l.Printf(fll.level, format, v...)
}

func (fll *fixedLevelLogger) Println(v ...interface{}) {
	fll.l.Println(fll.level, v...)
}

type tinyLogger struct {
	mu sync.RWMutex

	tags         map[string][]string
	destinations []*destination
}

func (tl *tinyLogger) SetLogLevel(level int, destinations ...Destination) {
	all := len(destinations) == 0

	ids := make(map[string]struct{})
	for _, dest := range destinations {
		ids[dest().ID()] = struct{}{}
	}

	tl.mu.Lock()

	for _, dest := range tl.destinations {
		if _, ok := ids[dest.ID()]; ok || all {
			dest.level = level
		}
	}

	tl.mu.Unlock()
}

func (tl *tinyLogger) GetFixedLevel(level int) FixedLevelLogger {
	return &fixedLevelLogger{tl, level}
}

func (tl *tinyLogger) AddTag(key string, value ...string) {
	tl.mu.Lock()
	tl.tags[key] = append(tl.tags[key], value...)
	tl.mu.Unlock()
}

func (tl *tinyLogger) Printf(level int, format string, v ...interface{}) {
	tl.output(level, fmt.Sprintf(format, v...), 1)
}

func (tl *tinyLogger) Println(level int, v ...interface{}) {
	tl.output(level, fmt.Sprint(v...), 1)
}

func (tl *tinyLogger) Fatalf(format string, v ...interface{}) {
	tl.output(Fatal, fmt.Sprintf(format, v...), 1)
	os.Exit(1)
}

func (tl *tinyLogger) Fatalln(v ...interface{}) {
	tl.output(Fatal, fmt.Sprint(v...), 1)
	os.Exit(1)
}

func (tl *tinyLogger) output(level int, message string, calldepth int) {
	tl.mu.RLock()
	for _, dest := range tl.destinations {
		if dest.level > level {
			continue
		}

		bytes := dest.formatter.GetOutput(level, message, tl.tags, calldepth+1)

		if _, err := dest.out.Write(bytes); err != nil {
			fmt.Printf(
				formatters.PaintText(
					formatters.ANSIColorRed,
					fmt.Sprintf("failed to write log to destination %s: %s",
						dest.ID(), err)))
		}
	}
	tl.mu.RUnlock()
}
