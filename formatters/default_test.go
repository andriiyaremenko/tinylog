package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	short       string = "hello logger"
	medium      string = "hello logger, my old friend"
	withNewLine string = "hello logger\nmy old friend"
	long        string = "hello logger, my old friend, lorem ipsum:\tlets break some lines:\tthis must be long enough to make it work:\twe need a little bit more text here"
)

const (
	lenWithFileWithoutTags int = 175 // length with colors
	lenDefault             int = 175 // length with colors
)

func TestDefaultLoggerFormatter(t *testing.T) {
	t.Run("GetOutput returns plain text rows", testGetOutputReturnsRows)
	t.Run("GetOutput returns plain text rows of exact length (175) if output has file but not tags",
		testGetOutputReturnsRowsOfExactLengthForWithFileWithoutTags)
	t.Run("GetOutput returns plain text rows of exact length (175)",
		testGetOutputReturnsRowsOfExactLength)
	t.Run("GetOutput can split long message into several rows",
		testGetOutputReturnsSeveralRowsForLongMessages)
	t.Run(`GetOutput will split message into several rows if there is \n in it`,
		testGetOutputReturnsSeveralRowsMessagesWithNewLines)
	t.Run("GetOutput would show file location for TRACE, DEBUG and FATAL",
		testGetOutputShowsFileForTraceDebugFatalOnly)
	t.Run("Test negative pad count case 1", testNegativePadCountCase_1)
}

func testGetOutputReturnsRows(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	b := f.GetOutput(0, short, make(map[string][]string), 0)
	s := string(b)
	r := []rune(s)

	assert.Equalf('\n', r[len(r)-1], "should end with new line: %s", s)
}

func testGetOutputReturnsRowsOfExactLengthForWithFileWithoutTags(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	b := f.GetOutput(0, PaintText(ANSIColorBlue, short), make(map[string][]string), 0)
	s := string(b)
	length := LenPrintableText(s[:len(s)-1])

	assert.Equal(lenWithFileWithoutTags, length, "should be of exact length")

	b = f.GetOutput(0, medium, make(map[string][]string), 0)
	s = string(b)
	length = LenPrintableText(s[:len(s)-1])

	assert.Equal(lenWithFileWithoutTags, length, "should be of exact length")

	b = f.GetOutput(0, long, make(map[string][]string), 0)
	s = string(b)

	for _, s := range strings.Split(s[:len(s)-1], "\n") {
		length = LenPrintableText(s)

		assert.Equal(lenWithFileWithoutTags, length, "should be of exact length")
	}
}

func testGetOutputReturnsRowsOfExactLength(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	tags := make(map[string][]string)
	tags["tag"] = []string{"cool tag"}
	b := f.GetOutput(0, PaintText(ColorFatal, short), tags, 0)
	s := string(b)
	length := LenPrintableText(s[:len(s)-1])

	assert.Equal(lenDefault, length, "should be of exact length")

	b = f.GetOutput(0, medium, tags, 0)
	s = string(b)
	length = LenPrintableText(s[:len(s)-1])

	assert.Equal(lenDefault, length, "should be of exact length")

	b = f.GetOutput(0, long, tags, 0)
	s = string(b)

	for _, s := range strings.Split(s[:len(s)-1], "\n") {
		length = LenPrintableText(s)

		assert.Equal(lenDefault, length, "should be of exact length")
	}
}

func testGetOutputReturnsSeveralRowsForLongMessages(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	b := f.GetOutput(0, long, make(map[string][]string), 0)
	s := string(b)
	// remove last new line
	rows := strings.Split(s[:len(s)-1], "\n")

	assert.Greater(len(rows), 1)
}

func testGetOutputReturnsSeveralRowsMessagesWithNewLines(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	b := f.GetOutput(0, withNewLine, make(map[string][]string), 0)
	s := string(b)
	// remove last new line
	rows := strings.Split(s[:len(s)-1], "\n")

	assert.Greater(len(rows), 1)
}

func testGetOutputShowsFileForTraceDebugFatalOnly(t *testing.T) {
	assert := assert.New(t)
	f := Default()
	for i := 0; i <= 5; i++ {
		b := f.GetOutput(i, short, make(map[string][]string), 0)
		if i <= 1 || i == 5 {
			// default_test.go - name of this file, where it was called
			assert.Containsf(string(b), "default_test.go", "should print file location for %d", i)
			continue
		}

		// default_test.go - name of this file, where it was called
		assert.NotContainsf(string(b), "default_test.go", "should not print file location for %d", i)
	}
}

func testNegativePadCountCase_1(t *testing.T) {
	message := "http://www.thessaliaradio.online/ - \x1b[32mattempt #1\x1b[0m \x1b[34m600ms\x1b[0m"
	assert := assert.New(t)
	f := Default()
	b := f.GetOutput(2, message, map[string][]string{
		"link": []string{
			`{"login":"test","password":"test123","cms":"","method":"","address":"http://perdetata.bg/","dateAdded":"2021-06-27T17:19:19.806708+03:00"}`,
		},
	}, 0)
	s := string(b)

	for _, s := range strings.Split(s[:len(s)-1], "\n") {
		println(s)
		length := LenPrintableText(s)

		assert.Equal(lenDefault, length, "should be of exact length")
	}
}
