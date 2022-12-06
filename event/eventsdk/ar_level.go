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

func (l level) Valid() bool {
	switch l.Self() {
	case "ERROR":
		return true
	case "WARN":
		return true
	case "INFO":
		return true
	default:
		return false
	}
}

func (l level) private() {
	// private 禁止用户自己实现接口。
}

const (
	ERROR level = "ERROR"
	WARN  level = "WARN"
	INFO  level = "INFO"
)
