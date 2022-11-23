package common

// ARLevel 定义了 Level 的3种级别。
type ARLevel interface {
	ERROR() Level
	WARN() Level
	INFO() Level

	private()
}

// Level 实际为 string 类型。
type Level string

// ERROR 返回 ERROR 级别的 Level 。
func (l Level) ERROR() Level {
	return "ERROR"
}

// WARN 返回 WARN 级别的 Level 。
func (l Level) WARN() Level {
	return "WARN"
}

// INFO 返回 INFO 级别的 Level 。
func (l Level) INFO() Level {
	return "INFO"
}

// private 禁止实现 ARLevel 接口。
func (l Level) private() {}

const (
	ERROR Level = "ERROR"
	WARN  Level = "WARN"
	INFO  Level = "INFO"
)
