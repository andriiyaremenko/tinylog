package tinylog

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/andriiyaremenko/tinylog/formatters"
)

// Returns new instance of `Logger` based on `out` and `formatter`
func NewLogger(out io.Writer, formatter formatters.LogFormatter) Logger {
	return &tinyLogger{
		out:       out,
		formatter: formatter,
		logLevel:  Info,
		tags:      make(map[string][]string),
	}
}

// Returns new instance of `Logger` based on `os.Stderr` as `out` and `formatters.Default()` as `formatter`
func NewDefaultLogger() Logger {
	return NewLogger(os.Stderr, formatters.Default())
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
	mu        sync.RWMutex
	out       io.Writer
	formatter formatters.LogFormatter
	logLevel  int
	tags      map[string][]string
}

func (tl *tinyLogger) SetLogLevel(level int) {
	tl.mu.Lock()
	tl.logLevel = level
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
	tl.print(level, fmt.Sprintf(format, v...))
}

func (tl *tinyLogger) Println(level int, v ...interface{}) {
	tl.print(level, fmt.Sprint(v...))
}

func (tl *tinyLogger) Fatalf(format string, v ...interface{}) {
	tl.print(Fatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (tl *tinyLogger) Fatalln(v ...interface{}) {
	tl.print(Fatal, fmt.Sprint(v...))
	os.Exit(1)
}

func (tl *tinyLogger) print(level int, message string) {
	tl.mu.RLock()

	if tl.logLevel > level {
		tl.mu.RUnlock()
		return
	}

	tl.mu.RUnlock()
	tl.output(level, message, 2)
}

func (tl *tinyLogger) output(level int, message string, calldepth int) {
	tl.mu.RLock()
	bytes := tl.formatter.GetOutput(level, message, tl.tags, calldepth+1)
	tl.mu.RUnlock()

	if _, err := tl.out.Write(bytes); err != nil {
		fmt.Printf(formatters.PaintText(formatters.ANSIColorRed, fmt.Sprintf("failed to write log to io.Writer: %s", err)))
	}
}
