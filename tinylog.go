package tinylog

import (
	"context"
	"fmt"
	"regexp"
)

const (
	// Most verbose level.
	// Is supposed to carry useful information to developers.
	// Is supposed to contain file and row number of place logging function was called.
	Trace int = iota
	// Is supposed to carry useful information not only to developers but to support as well.
	// Is supposed to contain file and row number of place logging function was called.
	Debug
	// Is supposed to carry useful information about program usage.
	Info
	// Is supposed to carry information about obscure behaviour of a program.
	Warn
	// Is supposed to carry information about error that caused abort of current operation.
	Error
	// Is supposed to carry information about error that caused program to exit.
	// Is supposed to contain file and row number of place logging function was called.
	Fatal
)

var (
	regexpTrace = regexp.MustCompile("(?i)trace")
	regexpDebug = regexp.MustCompile("(?i)debug")
	regexpInfo  = regexp.MustCompile("(?i)info")
	regexpWarn  = regexp.MustCompile("(?i)warn")
	regexpError = regexp.MustCompile("(?i)error")
	regexpFatal = regexp.MustCompile("(?i)fatal")
)

// Parses string to LogLevel int.
// Returns int if level is recognized or error if not.
func ParseLogLevel(level string) (int, error) {
	switch {
	case regexpTrace.MatchString(level):
		return Trace, nil
	case regexpDebug.MatchString(level):
		return Debug, nil
	case regexpInfo.MatchString(level):
		return Info, nil
	case regexpWarn.MatchString(level):
		return Warn, nil
	case regexpError.MatchString(level):
		return Error, nil
	case regexpFatal.MatchString(level):
		return Fatal, nil
	default:
		return 0, fmt.Errorf("unrecognized log level: %s", level)
	}
}

// parses string to LogLevel int.
func MustParseLogLevel(level string) int {
	l, err := ParseLogLevel(level)
	if err != nil {
		panic(err)
	}

	return l
}

// Destination configuration.
// Defines log formatter and log level for particular output.
type Destination func() *destination

type LogLevelSetter interface {
	// Sets verbosity level for set of Destinations.
	// 0 = Trace;
	// 1 = Debug;
	// 2 = Info;
	// 3 = Warn;
	// 4 = Error;
	// 5 = Fatal;
	// If no Destination were provided default LogLevelSetter Destinations are expected to be used.
	SetLogLevel(level int, destinations ...Destination)
}

// Logger bound to concrete log level.
type FixedLevelLogger interface {
	// Printf formats according to a format specifier and writes to io.Writer with appended newline.
	Printf(format string, v ...interface{})
	// Fprintln formats using the default formats for its operands and writes to io.Writer.
	// Spaces are always added between operands and a newline is appended.
	Println(v ...interface{})
}

// Logger can print log of different verbosity level.
type Logger interface {
	LogLevelSetter

	// Returns instance of FixedLevelLogger that shares tags with Logger instance.
	GetFixedLevel(level int) FixedLevelLogger
	// Adds tag to a logger and all instances of FixedLevelLogger created from this Logger.
	AddTag(key string, value ...string)

	// Printf formats according to a format specifier and writes to io.Writer with level of verbosity.
	// 0 = Trace;
	// 1 = Debug;
	// 2 = Info;
	// 3 = Warn;
	// 4 = Error;
	// 5 = Fatal;
	Printf(level int, format string, v ...interface{})
	// Fprintln formats using the default formats for its operands and writes to io.Writer with level of verbosity.
	// Spaces are always added between operands and a newline is appended.
	// 0 = Trace;
	// 1 = Debug;
	// 2 = Info;
	// 3 = Warn;
	// 4 = Error;
	// 5 = Fatal;
	Println(level int, v ...interface{})
	// Fatalf is equivalent to l.Printf(tinylog.Fatal) followed by a call to os.Exit(1).
	Fatalf(format string, v ...interface{})
	// Fatalln is equivalent to l.Println(tinylog.Fatal) followed by a call to os.Exit(1).
	Fatalln(v ...interface{})
}

// LoggerFactory manages Loggers instances verbosity levels and can get Logger instance bound to context.
type LoggerFactory interface {
	LogLevelSetter
	// Returns instance of Logger bound to provided ctx with listed Destinations.
	// If no Destination were provided default LoggerFactory Destinations are expected to be used.
	GetLogger(ctx context.Context, destinations ...Destination) Logger
	// Returns all Destinations for this LoggerFactory.
	Destinations() []Destination
}
