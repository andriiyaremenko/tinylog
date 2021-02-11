package tinylog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	t.Run("Test Trace level", getLevelTest("trace", "TRACE", "Trace", Trace))
	t.Run("Test Debug Level", getLevelTest("debug", "DEBUG", "Debug", Debug))
	t.Run("Test Info Level", getLevelTest("info", "INFO", "Info", Info))
	t.Run("Test Warn Level", getLevelTest("warn", "WARN", "Warn", Warn))
	t.Run("Test Error Level", getLevelTest("error", "ERROR", "Error", Error))
	t.Run("Test Fatal Level", getLevelTest("fatal", "FATAL", "Fatal", Fatal))
}

func getLevelTest(lowerCase, upperCase, capitalizedCase string, logLevel int) func(t *testing.T) {
	return func(t *testing.T) {
		assert := assert.New(t)

		level, err := ParseLogLevel(lowerCase)
		assert.NoError(err, "no error should be returned")
		assert.Equalf(logLevel, level, "lowercase level %s should be resolved correctly", lowerCase)

		level, err = ParseLogLevel(upperCase)
		assert.NoError(err, "no error should be returned")
		assert.Equalf(logLevel, level, "upper level %s should be resolved correctly", upperCase)

		level, err = ParseLogLevel(capitalizedCase)
		assert.NoError(err, "no error should be returned")
		assert.Equalf(logLevel, level, "capitalized level %s should be resolved correctly", capitalizedCase)
	}
}
