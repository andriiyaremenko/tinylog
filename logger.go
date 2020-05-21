package tinylog

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	ZeroModule = ""
)

func NewTinyLogger(out io.Writer, module string) TinyLogger {
	mPrefix := ""
	if module != "" {
		mPrefix = fmt.Sprintf("[%s]\t", module)
	}
	return &tinyLogger{
		debug:    log.New(out, fmt.Sprintf("[debug]\t%s", mPrefix), log.Ldate|log.Ltime|log.Lshortfile),
		info:     log.New(out, fmt.Sprintf("[info]\t%s", mPrefix), log.Ldate|log.Ltime|log.Lshortfile),
		warn:     log.New(out, fmt.Sprintf("[warn]\t%s", mPrefix), log.Ldate|log.Ltime|log.Lshortfile),
		err:      log.New(out, fmt.Sprintf("[err]\t%s", mPrefix), log.Ldate|log.Ltime|log.Lshortfile),
		fatal:    log.New(out, fmt.Sprintf("[fatal]\t%s", mPrefix), log.Ldate|log.Ltime|log.Lshortfile),
		logLevel: Info,
	}
}

func NewConsoleTinyLogger(module string) TinyLogger {
	return NewTinyLogger(os.Stderr, module)
}

type tinyLogger struct {
	debug    *log.Logger
	info     *log.Logger
	warn     *log.Logger
	err      *log.Logger
	fatal    *log.Logger
	logLevel int
}

func (tl *tinyLogger) Debug(v ...interface{}) {
	if tl.logLevel > Debug {
		return
	}
	tl.debug.Output(2, fmt.Sprintln(v...))
}

func (tl *tinyLogger) Debugf(format string, v ...interface{}) {
	if tl.logLevel > Debug {
		return
	}
	tl.debug.Output(2, fmt.Sprintf(format, v...))
}

func (tl *tinyLogger) Info(v ...interface{}) {
	if tl.logLevel > Info {
		return
	}
	tl.info.Output(2, fmt.Sprintln(v...))
}

func (tl *tinyLogger) Infof(format string, v ...interface{}) {
	if tl.logLevel > Info {
		return
	}
	tl.info.Output(2, fmt.Sprintf(format, v...))
}

func (tl *tinyLogger) Warn(v ...interface{}) {
	if tl.logLevel > Warn {
		return
	}
	tl.warn.Output(2, fmt.Sprintln(v...))
}

func (tl *tinyLogger) Warnf(format string, v ...interface{}) {
	if tl.logLevel > Warn {
		return
	}
	tl.warn.Output(2, fmt.Sprintf(format, v...))
}

func (tl *tinyLogger) Err(v ...interface{}) {
	if tl.logLevel > Err {
		return
	}
	tl.err.Output(2, fmt.Sprintln(v...))
}

func (tl *tinyLogger) Errf(format string, v ...interface{}) {
	if tl.logLevel > Err {
		return
	}
	tl.err.Output(2, fmt.Sprintf(format, v...))
}

func (tl *tinyLogger) Fatal(v ...interface{}) {
	tl.fatal.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (tl *tinyLogger) Fatalf(format string, v ...interface{}) {
	tl.fatal.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (tl *tinyLogger) SetLogLevel(level int) {
	tl.logLevel = level
}
