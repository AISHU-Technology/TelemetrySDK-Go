package common

// ARAttribute 对外暴露的 attribute 接口。
type ARAttribute interface {
	// Valid 校验 attribute 是否合法。
	Valid() bool
	// GetKey 返回 attribute 的键。
	GetKey() string
	// GetValue 返回 attribute 的值。
	GetValue() ARValue

	// private 禁止自己实现接口
	private()
}
