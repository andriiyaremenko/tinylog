package tinylog

import "context"

const (
	// Most verbose level
	// Is supposed to carry useful information to developers
	// Is supposed to contain file and row number of place logging function was called
	Trace int = iota
	// Is supposed to carry useful information not only to developers but to support as well
	// Is supposed to contain file and row number of place logging function was called
	Debug
	// Is supposed to carry useful information about program usage
	Info
	// Is supposed to carry information about obscure behaviour of a program
	Warn
	// Is supposed to carry information about error that caused abort of current operation
	Error
	// Is supposed to carry information about error that caused program to exit
	// Is supposed to contain file and row number of place logging function was called
	Fatal
)

type LogLevelSetter interface {
	// Sets verbosity level
	// `0` = `Trace`
	// `1` = `Debug`
	// `2` = `Info`
	// `3` = `Warn`
	// `4` = `Error`
	// `5` = `Fatal`
	SetLogLevel(level int)
}

// Logger bound to concrete log level
type FixedLevelLogger interface {
	// Printf formats according to a format specifier and writes to `io.Writer`
	Printf(format string, v ...interface{})
	// Fprintln formats using the default formats for its operands and writes to `io.Writer`
	// Spaces are always added between operands and a newline is appended
	Println(v ...interface{})
}

// Logger can print log of different verbosity level
type Logger interface {
	LogLevelSetter

	// Returns instance of `FixedLevelLogger` that shares tags with `Logger` instance
	GetFixedLevel(level int) FixedLevelLogger
	// Adds tag to a logger and all instances of `FixedLevelLogger` created from this `Logger`
	AddTag(key string, value ...string)

	// Printf formats according to a format specifier and writes to `io.Writer` with `level` of verbosity
	// `0` = `Trace`
	// `1` = `Debug`
	// `2` = `Info`
	// `3` = `Warn`
	// `4` = `Error`
	// `5` = `Fatal`
	Printf(level int, format string, v ...interface{})
	// Fprintln formats using the default formats for its operands and writes to `io.Writer` with `level` of verbosity
	// Spaces are always added between operands and a newline is appended
	// `0` = `Trace`
	// `1` = `Debug`
	// `2` = `Info`
	// `3` = `Warn`
	// `4` = `Error`
	// `5` = `Fatal`
	Println(level int, v ...interface{})
	// Fatalf is equivalent to `l.Printf(tinylog.Fatal)` followed by a call to `os.Exit(1)`
	Fatalf(format string, v ...interface{})
	// Fatalln is equivalent to `l.Println(tinylog.Fatal)` followed by a call to `os.Exit(1)`
	Fatalln(v ...interface{})
}

// `LoggerFactory` manages `Logger`s instances verbosity levels and can get `Logger` instance bound to context
type LoggerFactory interface {
	LogLevelSetter
	// Returns instance of `Logger` bound to provided `ctx`
	GetLogger(ctx context.Context) Logger
}
