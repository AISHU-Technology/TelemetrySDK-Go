package common

// ARLevel 定义了 level 的3种级别。
type ARLevel interface {
	// Self 返回事件级别。
	Self() string

	// private 禁止自己实现接口
	private()
	// ERROR 事件级别：错误。
	// WARN 事件级别：警告。
	// INFO 事件级别：一般。
}
