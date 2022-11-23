package common

import "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"

// Level 实际为 string 类型。
type Level string

// ERROR 返回 ERROR 级别的 Level 。
func (l Level) ERROR() model.ARLevel {
	return Level("ERROR")
}

// WARN 返回 WARN 级别的 Level 。
func (l Level) WARN() model.ARLevel {
	return Level("WARN")
}

// INFO 返回 INFO 级别的 Level 。
func (l Level) INFO() model.ARLevel {
	return Level("INFO")
}

// private 禁止实现 ARLevel 接口。
//func (l Level) private() {}

const (
	ERROR Level = "ERROR"
	WARN  Level = "WARN"
	INFO  Level = "INFO"
)
