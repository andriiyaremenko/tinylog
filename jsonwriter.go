package tinylog

import (
	"encoding/json"
	"io"
	"regexp"
	"time"
)

var (
	levelRegex     = regexp.MustCompile(`^[A-Z]+`)
	moduleRegex    = regexp.MustCompile(`\|\w+\|`)
	timeStampRegex = regexp.MustCompile(`\d{4}\/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`)
	locationRegex  = regexp.MustCompile(`\w+\.go:\d+`)
)

func NewJSONWriter(w io.Writer) io.Writer {
	return &jsonWriter{w}
}

type record struct {
	Level     string `json:"level"`
	Module    string `json:"module"`
	Timestamp int64  `json:"timestamp"`
	Location  string `json:"location"`
	Message   string `json:"message"`
}

type jsonWriter struct {
	w io.Writer
}

func (jw *jsonWriter) Write(p []byte) (n int, err error) {
	layout := "2006/01/02 15:04:05"
	level := levelRegex.Find(p)
	mod := moduleRegex.Find(p)
	t := timeStampRegex.Find(p)
	loc := locationRegex.Find(p)
	modLen := len(mod)
	add := 8 // length of white spaces and ":"
	if modLen == 0 {
		add--
	}
	mIndex := len(level) + modLen + len(t) + len(loc) + add
	mes := p[mIndex : len(p)-1]
	ts, err := time.Parse(layout, string(t))
	if err != nil {
		return
	}
	result, err := json.Marshal(&record{
		string(level),
		string(mod),
		ts.Unix(),
		string(loc),
		string(mes),
	})
	if err != nil {
		return
	}
	return jw.w.Write(result)
}
