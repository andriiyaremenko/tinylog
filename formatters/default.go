package formatters

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	totalSpace   int = 175
	messageSpace int = 50
)

// Returns default instance of Formatter with time.RFC822 as timeFormat.
// Formatter that returns log message in form of colorized plain text rows of fixed length.
func Default() LogFormatter {
	return New(time.RFC822)
}

// Returns default instance of Formatter with timeFormat.
// Formatter that returns log message in form of colorized plain text rows of fixed length.
func New(timeFormat string) LogFormatter {
	return &defaultFormatter{timeFormat: timeFormat}
}

type defaultFormatter struct {
	timeFormat string
}

func (df *defaultFormatter) GetOutput(level int, message string, tags map[string][]string, calldepth int) []byte {
	now := time.Now() // get this early.
	printFile := level <= 1 || level == 5

	if level == 0 || level == 5 {
		message = DecolorizeString(message)
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 2, 2, ' ', 0)

	fmt.Fprint(w, message)
	w.Flush()

	message = buf.String()

	levelS, color := getLevelTextAndColor(level)
	file, line := getFileAndLine(calldepth + 1)
	dateString := now.Format(df.timeFormat)
	foundAtString := fmt.Sprintf("at %v:%d", file, line)

	levelSection := PaintText(color, levelS) + " "
	dateSection := PaintText(ANSIColorGray, dateString) + " "
	fileSection := " " + PaintText(ANSIColorGray, foundAtString)

	spaceForLevel := len(levelS) + 1
	spaceForDate := len(dateString) + 1
	spaceForFile := LenPrintableText(foundAtString) + 1
	spaceForTags := 0
	messageLength := LenPrintableText(message)

	if !printFile {
		fileSection = ""
		spaceForFile = 0
	}

	var keys []string

	for k := range tags {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var tagsSectionBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			tagsSectionBuilder.WriteByte(';')
			tagsSectionBuilder.WriteByte(' ')

			spaceForTags += 2
		}

		tagsSectionBuilder.WriteString(PaintText(color, k))
		tagsSectionBuilder.WriteByte('=')

		spaceForTags += LenPrintableText(k) + 1

		for j, s := range tags[k] {
			if j > 0 {
				tagsSectionBuilder.WriteByte(',')

				spaceForTags++
			}

			tagsSectionBuilder.WriteString(PaintText(color, s))

			spaceForTags += LenPrintableText(s)
		}
	}

	tagsSection := " " + tagsSectionBuilder.String()

	if len(tagsSection) > messageSpace {
		panic(
			fmt.Sprintf(
				"tag %q is too long - consider moving part of into message or using other formatter",
				strings.Trim(DecolorizeString(tagsSection), " "),
			),
		)
	}

	spaceForTags++

	spaceForEverythingElse := spaceForLevel + spaceForDate + spaceForFile + spaceForTags
	spaceForMessage := totalSpace - spaceForEverythingElse
	messageSpace := messageSpace

	if messageLength < messageSpace {
		messageSpace = messageLength
	}

	fileOnSeparateLine := false
	tagsAndFileOnSeparateLine := false

	for {
		if !fileOnSeparateLine && spaceForMessage < messageSpace {
			fileOnSeparateLine = true
			spaceForMessage += spaceForFile
			continue
		}

		if spaceForMessage < messageSpace {
			tagsAndFileOnSeparateLine = true
			spaceForMessage += spaceForTags
		}

		break
	}

	for strings.HasSuffix(message, "\n") {
		message = strings.TrimSuffix(message, "\n")
	}

	var b []byte
	if messageLength > spaceForMessage || strings.Contains(message, "\n") {
		for _, messagePart := range splitMessageIntoRows(message, spaceForMessage) {
			spaceForMessagePart := LenPrintableText(messagePart)

			if level == 0 {
				messagePart = PaintText(ColorTrace, messagePart)
			}

			if level == 5 {
				messagePart = PaintText(ColorFatal, messagePart)
			}

			b = append(b, getFinalOutput(levelSection, spaceForLevel, dateSection, spaceForDate,
				tagsSection, spaceForTags, tagsAndFileOnSeparateLine,
				fileSection, spaceForFile, fileOnSeparateLine,
				messagePart, spaceForMessagePart, totalSpace)...)
			b = append(b, '\n')
		}

		return b
	}

	if level == 0 {
		message = PaintText(ColorTrace, message)
	}

	if level == 5 {
		message = PaintText(ColorFatal, message)
	}

	b = append(b, getFinalOutput(levelSection, spaceForLevel, dateSection, spaceForDate,
		tagsSection, spaceForTags, tagsAndFileOnSeparateLine,
		fileSection, spaceForFile, fileOnSeparateLine,
		message, messageLength, totalSpace)...)
	b = append(b, '\n')

	return b
}

func getFinalOutput(levelSection string, levelLength int, dateSection string, dateLength int,
	tagsSection string, tagsLength int, tagsAndFileAreOnSeparateLine bool,
	fileSection string, fileLength int, fileIsOnSeparateLine bool,
	message string, messageLength int, totalLength int) []byte {
	var b []byte
	b = append(b, []byte(levelSection)...)
	b = append(b, []byte(dateSection)...)
	b = append(b, []byte(message)...)
	padSpaces := func(count int) []byte {
		if count < 0 {
			switch {
			case tagsLength > 0 && fileLength > 0:
				panic(
					fmt.Sprintf("message %q\n\twith tags %q\n\tand file %q\n\tcaused negative (%d) pad spaces count",
						message,
						strings.Trim(DecolorizeString(tagsSection), " "),
						strings.Trim(DecolorizeString(fileSection), " "),
						count,
					),
				)
			case tagsLength > 0:
				panic(
					fmt.Sprintf("message %q\n\twith tags %q\n\tcaused negative (%d) pad spaces count",
						message,
						strings.Trim(DecolorizeString(tagsSection), " "),
						count,
					),
				)
			case fileLength > 0:
				panic(
					fmt.Sprintf("message %q with file %q\n\tcaused negative (%d) pad spaces count",
						message,
						strings.Trim(DecolorizeString(fileSection), " "),
						count,
					),
				)
			default:
				panic(fmt.Sprintf("message %q caused negative (%d) pad spaces count", message, count))
			}
		}

		return []byte(strings.Repeat(" ", count))
	}

	switch {
	case tagsAndFileAreOnSeparateLine:
		b = append(b, padSpaces(totalLength-levelLength-dateLength-messageLength)...)
		b = append(b, '\n')
		b = append(b, []byte(levelSection)...)
		b = append(b, []byte(dateSection)...)
		b = append(b, padSpaces(totalLength-levelLength-dateLength-tagsLength-fileLength)...)
		b = append(b, []byte(fileSection)...)
		b = append(b, []byte(tagsSection)...)
	case fileIsOnSeparateLine:
		b = append(b, padSpaces(totalLength-levelLength-dateLength-messageLength-tagsLength)...)
		b = append(b, []byte(tagsSection)...)
		b = append(b, '\n')
		b = append(b, []byte(levelSection)...)
		b = append(b, []byte(dateSection)...)
		b = append(b, padSpaces(totalLength-levelLength-dateLength-fileLength)...)
		b = append(b, []byte(fileSection)...)
	default:
		b = append(b, padSpaces(totalLength-levelLength-dateLength-messageLength-tagsLength-fileLength)...)
		b = append(b, []byte(fileSection)...)
		b = append(b, []byte(tagsSection)...)
	}

	return b
}

func splitMessageIntoRows(message string, spaceForMessage int) []string {
	parts := make([]string, 0, 1)
	messageRows := strings.Split(message, "\n")

	for _, messageRow := range messageRows {
		if len(messageRow) == 0 {
			continue
		}

		if LenPrintableText(messageRow) > spaceForMessage {
			parts = append(parts,
				splitMessageBySignsFunc(";",
					splitMessageBySignsFunc(",", splitMessageIntoParts))(messageRow, spaceForMessage)...)
			continue
		}

		parts = append(parts, messageRow)
	}

	return parts
}

func splitMessageBySignsFunc(sign string, nextSplitFunc func(string, int) []string) func(string, int) []string {
	return func(message string, spaceForMessage int) []string {
		parts := make([]string, 0, 1)
		messageRows := strings.Split(message, sign)

		for i, messageRow := range messageRows {
			if len(messageRow) == 0 {
				continue
			}

			if strings.HasPrefix(message, sign) && i == 0 {
				messageRow = sign + messageRow
			}

			if len(messageRows) > 1 && i < len(messageRows)-1 {
				messageRow = messageRow + sign
			}

			if LenPrintableText(messageRow) > spaceForMessage {
				parts = append(parts, nextSplitFunc(messageRow, spaceForMessage)...)
				continue
			}

			parts = append(parts, messageRow)
		}

		return parts
	}
}

func splitMessageIntoParts(messageRow string, spaceForMessage int) []string {
	parts := make([]string, 0, 1)
	messageParts := strings.SplitN(messageRow, ": ", 2)

	for i, messagePart := range messageParts {
		if len(messagePart) == 0 {
			continue
		}

		if i < len(messageParts)-1 {
			messagePart = messagePart + ":"
		}

		if i > 0 {
			messagePart = " " + messagePart
		}

		isTooLong := LenPrintableText(messagePart) > spaceForMessage
		if isTooLong && strings.Contains(messagePart, ": ") {
			parts = append(parts, splitMessageIntoParts(messagePart, spaceForMessage)...)
			continue
		}

		if isTooLong {
			parts = append(parts, breakMessageInLines(messagePart, spaceForMessage)...)
			continue
		}

		parts = append(parts, messagePart)
	}

	return parts
}

func breakMessageInLines(messagePart string, spaceForMessage int) []string {
	messageLength := len(messagePart)
	nParts := int(math.Ceil(float64(messageLength) / float64(spaceForMessage)))
	parts := make([]string, 0, 1)

	for i := 0; i < nParts; i++ {
		startShift := spaceForMessage
		if i > 0 {
			startShift -= 1
		}

		start := i * startShift
		finish := start + startShift

		if finish > messageLength {
			finish = messageLength
		}

		nextPart := messagePart[start:finish]

		if len(nextPart) == 0 {
			continue
		}

		if i > 0 {
			nextPart = " " + nextPart
		}

		parts = append(parts, nextPart)

		if finish == messageLength {
			break
		}
	}

	return parts
}
