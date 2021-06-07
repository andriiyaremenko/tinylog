package tinylog

import (
	"bytes"
	"testing"

	"github.com/andriiyaremenko/tinylog/formatters"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("Default log level is Info", testDefaultLogLevelIsInfo)
	t.Run("SetLogLevel changes verbosity level", testSetLogLevel)
	t.Run("SetLogLevel changes verbosity level for particular Destination", testSetLogLevelForDestination)
	t.Run("AddTag adds tag to output", testAddTag)
	t.Run("GetFixedLevel returns FixedLevelLogger of correct level", testGetFixedLevel)
	t.Run("FixedLevelLogger respects verbosity level", testFixedLevelRespectsVerbosity)
}

func testDefaultLogLevelIsInfo(t *testing.T) {
	assert := assert.New(t)
	l, _, b := getLogger()

	l.Println(Trace, "trace")
	l.Println(Debug, "debug")
	l.Println(Info, "info")

	result := b.String()
	assert.Contains(result, "info", "default log level should be Info")
	assert.NotContains(result, "trace", "default log level should be Info")
	assert.NotContains(result, "debug", "default log level should be Info")
}

func testSetLogLevel(t *testing.T) {
	assert := assert.New(t)
	l, destinations, b := getLogger()

	l.SetLogLevel(Debug, destinations...)
	l.Println(Trace, "trace")
	l.Println(Debug, "debug")

	result := b.String()
	assert.Contains(result, "debug", "SetLogLevel(Debug): log level should be Debug")
	assert.NotContains(result, "trace", "SetLogLevel(Debug): log level should be Debug")

	l, _, b = getLogger()

	l.SetLogLevel(Trace)
	l.Println(Trace, "trace")

	result = b.String()

	assert.Contains(result, "trace", "SetLogLevel(Trace): log level should be Trace")

	l, _, b = getLogger()

	l.SetLogLevel(Error)
	l.Println(Trace, "trace")
	l.Println(Debug, "debug")
	l.Println(Info, "info")
	l.Println(Warn, "warn")
	l.Println(Error, "error")

	result = b.String()

	assert.Contains(result, "error", "SetLogLevel(Error): log level should be Error")
	assert.NotContains(result, "trace", "SetLogLevel(Error): log level should be Error")
	assert.NotContains(result, "debug", "SetLogLevel(Error): log level should be Error")
	assert.NotContains(result, "info", "SetLogLevel(Error): log level should be Error")
	assert.NotContains(result, "warn", "SetLogLevel(Error): log level should be Error")
}

func testSetLogLevelForDestination(t *testing.T) {
	assert := assert.New(t)
	b1 := new(bytes.Buffer)
	destination1 := DestinationFunc(b1, formatters.Default(), Info)
	b2 := new(bytes.Buffer)
	destination2 := DestinationFunc(b2, formatters.Default(), Info)
	l := NewLogger(destination1, destination2)

	l.SetLogLevel(Debug, destination1)
	l.Println(Trace, "trace")
	l.Println(Debug, "debug")

	result := b1.String()
	assert.Contains(result, "debug", "SetLogLevel(Debug) for destination1: destination1 log level should be Debug")
	assert.NotContains(result, "trace", "SetLogLevel(Debug) for destination1: destination1 log level should be Debug")

	result = b2.String()
	assert.NotContains(result, "debug", "SetLogLevel(Debug) for destination1: destination2 log level should be Info")
	assert.NotContains(result, "trace", "SetLogLevel(Debug) for destination1: destination2 log level should be Info")
}

func testAddTag(t *testing.T) {
	assert := assert.New(t)
	l, _, b := getLogger()

	l.AddTag("user", "me", "cat")
	l.Println(Info, "info")

	result := b.String()

	assert.Contains(result, "info", "message should be printed")
	assert.Contains(result, "user", "tag should be printed")
	assert.Contains(result, "me", "tag should be printed")
	assert.Contains(result, "cat", "tag should be printed")
}

func testGetFixedLevel(t *testing.T) {
	assert := assert.New(t)
	l, _, b := getLogger()

	l.AddTag("user", "me", "cat")

	ll := l.GetFixedLevel(Warn)

	ll.Println("ooops")

	result := b.String()

	assert.Contains(result, "WARN", "verbosity level should be printed")
	assert.Contains(result, "ooops", "message should be printed")
	assert.Contains(result, "user", "tag still should be printed")
	assert.Contains(result, "me", "tag still should be printed")
	assert.Contains(result, "cat", "tag still should be printed")
}

func testFixedLevelRespectsVerbosity(t *testing.T) {
	assert := assert.New(t)
	l, _, b := getLogger()

	l.AddTag("user", "me", "cat")

	ll := l.GetFixedLevel(Debug)

	ll.Println("woof")

	result := b.String()

	assert.Equal(result, "", "nothing should be printed")

	l.SetLogLevel(Debug)
	ll.Println("woof")

	result = b.String()

	assert.Contains(result, "DEBUG", "verbosity level should be printed")
	assert.Contains(result, "woof", "message should be printed")
	assert.Contains(result, "user", "tag still should be printed")
	assert.Contains(result, "me", "tag still should be printed")
	assert.Contains(result, "cat", "tag still should be printed")
}
