package tinylog

type Record struct {
	Level     string `json:"level"`
	Module    string `json:"module"`
	Location  string `json:"location"`
	Message   string `json:"message"`
	TimeStamp int64  `json:"timeStamp"`
}
