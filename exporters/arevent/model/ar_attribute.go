package model

// ARAttribute 对外暴露的 Attribute 接口。
type ARAttribute interface {
	// Valid 校验 Attribute 是否合法。
	Valid() bool
	// GetKey 返回 Attribute 的键。
	GetKey() string
	// GetValue 返回 Attribute 的值。
	GetValue() ARValue
}

// ARValue 对外暴露的 Value 接口。
type ARValue interface {
	// GetType 返回 Value 的类型。
	GetType() string
	// GetValue 返回 Value 的值。
	GetValue() interface{}
}
