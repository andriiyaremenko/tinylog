package formatters

type Log struct {
	LevelCode     int                 `json:"levelCode"`
	Level         string              `json:"level"`
	Location      string              `json:"location"`
	Message       string              `json:"message"`
	Tags          map[string][]string `json:"tags"`
	TimeStampUnix int64               `json:"timeStamp"`
}
