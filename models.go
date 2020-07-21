package tinylog

type Record struct {
	LevelCode int                 `json:"levelCode"`
	Level     string              `json:"level"`
	Module    string              `json:"module"`
	Location  string              `json:"location"`
	Message   string              `json:"message"`
	TimeStamp int64               `json:"timeStamp"`
	Tags      map[string][]string `json:"tags"`
}
