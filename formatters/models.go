package formatters

import (
	_ "encoding/json"
	"time"
)

// Log model returned by JSONFormatter.
type Log struct {
	LevelCode int                 `json:"levelCode"`
	Level     string              `json:"level"`
	Location  string              `json:"location"`
	Message   string              `json:"message"`
	Tags      map[string][]string `json:"tags"`
	DateUnix  time.Time           `json:"date"`
}
