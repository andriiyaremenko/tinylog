package tinylog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

type format int

const (
	String format = iota
	JSON
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"
	bold        = "\033[1m"
)

const (
	NilModule = ""
)

func NewTinyLogger(out io.Writer, format format, module, timeFormat string) Logger {
	return &tinyLogger{
		out:        out,
		format:     format,
		module:     module,
		timeFormat: timeFormat,
		logLevel:   Info,
		tags:       make(map[string][]string),
	}
}

func NewConsoleTinyLogger(module, timeFormat string) Logger {
	return NewTinyLogger(os.Stderr, String, module, timeFormat)
}

type tinyLogger struct {
	mu         sync.RWMutex
	out        io.Writer
	format     format
	module     string
	timeFormat string
	logLevel   logLevel
	tags       map[string][]string
}

func (tl *tinyLogger) AddTag(ctx context.Context, key string, value ...string) {
	tl.mu.Lock()
	tl.tags[key] = value
	tl.mu.Unlock()

	go func() {
		<-ctx.Done()
		tl.mu.Lock()
		tl.tags = make(map[string][]string)
		tl.mu.Unlock()
	}()
}

func (tl *tinyLogger) Debug(v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Debug {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintln(v...), Debug)
}

func (tl *tinyLogger) Debugf(format string, v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Debug {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintf(format, v...), Debug)
}

func (tl *tinyLogger) Info(v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Info {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintln(v...), Info)
}

func (tl *tinyLogger) Infof(format string, v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Info {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintf(format, v...), Info)
}

func (tl *tinyLogger) Warn(v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Warn {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintln(v...), Warn)
}

func (tl *tinyLogger) Warnf(format string, v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Warn {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintf(format, v...), Warn)
}

func (tl *tinyLogger) Error(v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Error {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintln(v...), Error)
}

func (tl *tinyLogger) Errorf(format string, v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Error {
		return
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintf(format, v...), Error)
}

func (tl *tinyLogger) Fatal(v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Fatal {
		tl.mu.RUnlock()
		os.Exit(1)
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintln(v...), Fatal)
	os.Exit(1)
}

func (tl *tinyLogger) Fatalf(format string, v ...interface{}) {
	tl.mu.RLock()
	if tl.logLevel > Fatal {
		tl.mu.RUnlock()
		os.Exit(1)
	}
	tl.mu.RUnlock()
	tl.Output(2, fmt.Sprintf(format, v...), Fatal)
	os.Exit(1)
}

func (tl *tinyLogger) SetLogLevel(level logLevel) {
	tl.mu.Lock()
	tl.logLevel = level
	tl.mu.Unlock()
}

func (tl *tinyLogger) Output(calldepth int, message string, level logLevel) (err error) {
	now := time.Now() // get this early.
	var color string
	var levelS string
	var file string
	var line int
	var ok bool

	tl.mu.Lock()
	defer tl.mu.Unlock()

	switch level {
	case Debug:
		levelS = "DEBUG"
		color = colorGreen
	case Info:
		levelS = "INFO"
		color = colorCyan
	case Warn:
		levelS = "WARN"
		color = colorYellow
	case Error:
		levelS = "ERROR"
		color = colorRed
	case Fatal:
		levelS = "FATAL"
		color = bold + colorRed
	}

	tl.mu.Unlock()
	_, file, line, ok = runtime.Caller(calldepth)

	if !ok {
		file = "???"
		line = 0
	}
	tl.mu.Lock()

	var bytes []byte
	switch tl.format {
	case String:
		short := file

		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short

		bytes = append(bytes, colorGray...)
		bytes = append(bytes, '[')
		bytes = append(bytes, now.Format(tl.timeFormat)...)
		bytes = append(bytes, ']')
		bytes = append(bytes, colorReset...)
		bytes = append(bytes, ' ')
		bytes = append(bytes, color...)
		bytes = append(bytes, levelS...)
		bytes = append(bytes, '\t')

		if tl.module != "" {
			bytes = append(bytes, tl.module...)
			bytes = append(bytes, ' ')
		}

		bytes = append(bytes, colorReset...)
		bytes = append(bytes, colorGray...)
		bytes = append(bytes, fmt.Sprintf("at %v:%d", file, line)...)
		bytes = append(bytes, colorReset...)
		bytes = append(bytes, ' ')

		if len(tl.tags) > 0 {
			bytes = append(bytes, colorGray...)
			bytes = append(bytes, '{')
		}

		var keys []string

		for k := range tl.tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := tl.tags[k]
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, color...)
			bytes = append(bytes, k...)
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, colorGray...)
			bytes = append(bytes, ':')
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, color...)
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, colorGray...)
			bytes = append(bytes, '[')

			for _, s := range v {
				bytes = append(bytes, colorReset...)
				bytes = append(bytes, color...)
				bytes = append(bytes, s...)
				bytes = append(bytes, colorReset...)
				bytes = append(bytes, colorGray...)
				bytes = append(bytes, ',')
				bytes = append(bytes, ' ')
			}

			bytes = bytes[:len(bytes)-2]
			bytes = append(bytes, ']')
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, colorGray...)
			bytes = append(bytes, ',')
			bytes = append(bytes, ' ')
		}

		if len(tl.tags) > 0 {
			bytes = bytes[:len(bytes)-2]
			bytes = append(bytes, '}')
			bytes = append(bytes, colorReset...)
			bytes = append(bytes, ' ')
		}
		bytes = append(bytes, color...)
		bytes = append(bytes, message...)

		if len(message) == 0 || message[len(message)-1] != '\n' {
			bytes = append(bytes, '\n')
		}
		bytes = append(bytes, colorReset...)
	case JSON:
		if message[len(message)-1] == '\n' {
			message = message[:len(message)-1]
		}
		r := Record{
			LevelCode: int(level),
			Level:     levelS,
			Location:  fmt.Sprintf("%v:%d", file, line),
			Module:    tl.module,
			TimeStamp: now.Unix(),
			Message:   message,
			Tags:      tl.tags,
		}
		bytes, err = json.Marshal(r)

		if err != nil {
			return err
		}
	}
	_, err = tl.out.Write(bytes)
	return err
}
