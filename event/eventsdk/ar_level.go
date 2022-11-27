package eventsdk

// level 实际为 string 类型。
type level string

// newLevel 创建新的 level 。
func newLevel(l string) level {
	return level(l)
}

func (l level) Self() string {
	return string(l)
}

func (l level) private() {}

const (
	ERROR level = "ERROR"
	WARN  level = "WARN"
	INFO  level = "INFO"
)
