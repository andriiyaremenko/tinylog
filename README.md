# tinylog

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/.)

## Functions

### func [DefaultDestination](/destination.go#L22)

`func DefaultDestination() *destination`

Destination based on os.Stderr as out and formatters.Default() as formatter.

### func [MustParseLogLevel](/tinylog.go#L59)

`func MustParseLogLevel(level string) int`

parses string to LogLevel int.

### func [ParseLogLevel](/tinylog.go#L39)

`func ParseLogLevel(level string) (int, error)`

Parses string to LogLevel int.
Returns int if level is recognized or error if not.

## Types

### type [Destination](/tinylog.go#L70)

`type Destination func() *destination`

Destination configuration.
Defines log formatter and log level for particular output.

### type [FixedLevelLogger](/tinylog.go#L85)

`type FixedLevelLogger interface { ... }`

Logger bound to concrete log level.

### type [LogLevelSetter](/tinylog.go#L72)

`type LogLevelSetter interface { ... }`

### type [Logger](/tinylog.go#L94)

`type Logger interface { ... }`

Logger can print log of different verbosity level.

### type [LoggerFactory](/tinylog.go#L126)

`type LoggerFactory interface { ... }`

LoggerFactory manages Loggers instances verbosity levels and can get Logger instance bound to context.

## Sub Packages

* [formatters](./formatters)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
