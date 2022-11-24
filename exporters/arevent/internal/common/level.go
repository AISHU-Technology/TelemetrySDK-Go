package common

import "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"

// Level 实际为 string 类型。
type Level string

func (l Level) ERROR() model.ARLevel {
	return Level("ERROR")
}

func (l Level) WARN() model.ARLevel {
	return Level("WARN")
}

func (l Level) INFO() model.ARLevel {
	return Level("INFO")
}

const (
	ERROR Level = "ERROR"
	WARN  Level = "WARN"
	INFO  Level = "INFO"
)
