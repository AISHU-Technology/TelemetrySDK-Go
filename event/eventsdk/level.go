package eventsdk

// Level 定义了 level 的3种级别。
type Level interface {
	// Self 返回事件级别。
	Self() string
	// Valid 校验是否合法。
	Valid() bool
	// private 禁止用户自己实现接口。
	private()
	// ERROR 事件级别：错误。
	// WARN 事件级别：警告。
	// INFO 事件级别：一般。
}
