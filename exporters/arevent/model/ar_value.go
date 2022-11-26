package model

// ARValue 对外暴露的 Data 接口。
type ARValue interface {
	// GetType 返回 Data 的类型。
	GetType() string
	// GetData 返回 Data 的值。
	GetData() interface{}
	// private 禁止自己实现接口
	//private()
}
