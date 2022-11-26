package model

// ARAttribute 对外暴露的 Attribute 接口。
type ARAttribute interface {
	// Valid 校验 Attribute 是否合法。
	Valid() bool
	// GetKey 返回 Attribute 的键。
	GetKey() string
	// GetValue 返回 Attribute 的值。
	GetValue() ARValue
	// private 禁止自己实现接口
	//private()
}
