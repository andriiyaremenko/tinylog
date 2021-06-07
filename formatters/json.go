package formatters

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Formatter that returns log message in form of JSON.
const JSONFormatter jsonFormatter = "JSONFormatter"

type jsonFormatter string

func (f jsonFormatter) GetOutput(level int, message string, tags map[string][]string, calldepth int) []byte {
	now := time.Now().Round(time.Millisecond)
	levelS, _ := getLevelTextAndColor(level)
	file, line := getFileAndLine(calldepth + 1)
	message = DecolorizeString(message)

	m := Log{
		LevelCode: level,
		Level:     strings.TrimLeft(levelS, " "),
		Location:  fmt.Sprintf("%v:%d", file, line),
		Message:   message,
		DateUnix:  now,
		Tags:      tags}

	b, err := json.Marshal(m)

	if err != nil {
		fmt.Printf(PaintText(ANSIColorRed, fmt.Sprintf("%s: failed to write log: %s", f, err)))

		return []byte("")
	}

	return append(b, '\n')
}
