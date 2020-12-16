package formatters

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONLoggerFormatter(t *testing.T) {
	t.Run("GetOutput returns correct JSON of formatter.Log model", testJSONFormatterOutput)
}

func testJSONFormatterOutput(t *testing.T) {
	assert := assert.New(t)

	tags := make(map[string][]string)
	tags["tag"] = []string{"cool tag"}
	now := time.Now().Unix()
	b := JSONFormatter.GetOutput(2, "test json", tags, 0)
	m := new(Log)

	if err := json.Unmarshal(b, m); err != nil {
		assert.FailNow("got wrong log format")
	}

	expected := Log{
		LevelCode: 2,
		Level:     "INFO",
		Location:  "json_test.go:21",
		Message:   "test json",
		Tags:      tags,
		DateUnix:  now}

	assert.EqualValues(expected, *m, "log should contain all fields with correct values")
}
