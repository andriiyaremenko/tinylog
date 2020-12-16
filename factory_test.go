package tinylog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerFactory(t *testing.T) {
	t.Run("GetLogger returns Logger instance", testGetLogger)
	t.Run("GetLogger returns different Logger instance per Context", testGetLoggerDiffersPerContext)
	t.Run("GetLogger returns same Logger instance for same Context", testGetLoggerSameForOneContext)
	t.Run("SetLogLevel sets log level for all Logger instances", testLoggersRespectLogLevel)
}

func testGetLogger(t *testing.T) {
	assert := assert.New(t)
	lf, _ := getLoggerFactory()
	ctx := context.TODO()
	l := lf.GetLogger(ctx)

	assert.NotNil(l, "should return instance of Logger")
	assert.Implements((*Logger)(nil), l, "returned instance should implement Logger")
}

func testGetLoggerDiffersPerContext(t *testing.T) {
	assert := assert.New(t)
	lf, _ := getLoggerFactory()
	ctx1 := context.TODO()
	l1 := lf.GetLogger(ctx1)
	ctx2 := context.Background()
	l2 := lf.GetLogger(ctx2)
	l2.AddTag("tag", "serious")

	assert.NotNil(l1, "should return instance of Logger")
	assert.NotNil(l2, "should return instance of Logger")
	assert.NotEqual(l1, l2, "Loggers should be different for different contexts")
}

func testGetLoggerSameForOneContext(t *testing.T) {
	assert := assert.New(t)
	lf, _ := getLoggerFactory()
	ctx := context.TODO()
	l1 := lf.GetLogger(ctx)
	l2 := lf.GetLogger(ctx)

	assert.NotNil(l1, "should return instance of Logger")
	assert.NotNil(l2, "should return instance of Logger")
	assert.Equal(l1, l2, "Logger should be one for one context")
}

func testLoggersRespectLogLevel(t *testing.T) {
	assert := assert.New(t)
	lf, b := getLoggerFactory()
	ctx1 := context.TODO()
	l1 := lf.GetLogger(ctx1)
	ctx2 := context.TODO()
	l2 := lf.GetLogger(ctx2)
	ll1 := l1.GetFixedLevel(Debug)
	ll2 := l2.GetFixedLevel(Trace)

	lf.SetLogLevel(Error)

	l1.Println(Warn, "warn")
	l2.Println(Warn, "info")
	ll1.Println("debug")
	ll2.Println("trace")

	result := b.String()

	assert.Equal(result, "", "nothing should be printed")

	lf.SetLogLevel(Trace)

	l1.Println(Warn, "warn")
	l2.Println(Warn, "info")
	ll1.Println("debug")
	ll2.Println("trace")

	result = b.String()

	assert.Contains(result, "trace", "message should be printed")
	assert.Contains(result, "debug", "message should be printed")
	assert.Contains(result, "info", "message should be printed")
	assert.Contains(result, "warn", "message should be printed")
}
