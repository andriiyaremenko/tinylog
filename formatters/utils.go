package formatters

import (
	"regexp"
	"runtime"
	"unicode/utf8"
)

type Color string

const (
	ANSIReset       Color = "\033[0m"
	ANSIColorRed    Color = "\033[31m"
	ANSIColorGreen  Color = "\033[32m"
	ANSIColorYellow Color = "\033[33m"
	ANSIColorBlue   Color = "\033[34m"
	ANSIColorPurple Color = "\033[35m"
	ANSIColorCyan   Color = "\033[36m"
	ANSIColorWhite  Color = "\033[37m"
	ANSIColorGray   Color = "\033[90m"
	ANSIFontBold    Color = "\033[1m"

	ColorTrace Color = ANSIColorGray
	ColorDebug Color = ANSIColorCyan
	ColorInfo  Color = ANSIColorGreen
	ColorWarn  Color = ANSIColorYellow
	ColorError Color = ANSIColorRed
	ColorFatal Color = ANSIFontBold + ANSIColorRed
)

var ansiiColorMatch = regexp.MustCompile("\u001B\\[[;\\d]*m")

// Returns text prepended by ANSI color code and appended by ANSI color reset code.
func PaintText(color Color, text string) string {
	return string(color) + text + string(ANSIReset)
}

// Returns []byte prepended by ANSI color code and appended by ANSI color reset code.
func PaintBuffer(color Color, text []byte) []byte {
	var result []byte
	result = append(result, []byte(color)...)
	result = append(result, text...)
	result = append(result, []byte(ANSIReset)...)
	return result
}

// Returns text without color.
func DecolorizeString(text string) string {
	return ansiiColorMatch.ReplaceAllString(text, "")
}

// Returns text without color.
func DecolorizeBuffer(text []byte) []byte {
	return ansiiColorMatch.ReplaceAll(text, []byte(""))
}

// Returns utf8.RuneCountInString of text without ANSI color codes.
func LenPrintableText(text string) int {
	return utf8.RuneCountInString(string(ansiiColorMatch.ReplaceAll([]byte(text), []byte(""))))
}

func getLevelTextAndColor(level int) (string, Color) {
	var color Color
	var levelS string

	switch level {
	case 0:
		levelS = "TRACE"
		color = ColorTrace
	case 1:
		levelS = "DEBUG"
		color = ColorDebug
	case 2:
		levelS = " INFO"
		color = ColorInfo
	case 3:
		levelS = " WARN"
		color = ColorWarn
	case 4:
		levelS = "ERROR"
		color = ColorError
	case 5:
		levelS = "FATAL"
		color = ColorFatal
	}

	return levelS, color
}

func getFileAndLine(calldepth int) (string, int) {
	_, file, line, ok := runtime.Caller(calldepth + 1)

	if !ok {
		file = "???"
		line = 0
	}

	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	return short, line
}
