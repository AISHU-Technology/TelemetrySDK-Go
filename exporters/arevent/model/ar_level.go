package model

// ARLevel 定义了 Level 的3种级别。
type ARLevel interface {
	// Self 返回事件级别。
	Self() string
	// ERROR 事件级别：错误。
	ERROR() ARLevel
	// WARN 事件级别：警告。
	WARN() ARLevel
	// INFO 事件级别：一般。
	INFO() ARLevel
}
