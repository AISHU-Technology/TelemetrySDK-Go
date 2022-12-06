package eventsdk

// Attribute 对外暴露的 attribute 接口。
type Attribute interface {
	// Valid 校验 attribute 是否合法。
	Valid() bool
	// GetKey 返回 attribute 的键。
	GetKey() string
	// GetValue 返回 attribute 的值。
	GetValue() Value

	// private 禁止用户自己实现接口。
	private()
}
