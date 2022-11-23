package model

type ARAttribute interface {
	Valid() bool
	GetKey() string
	GetValue() ARValue

	//private()
}

type ARValue interface {
	GetType() string
	GetValue() interface{}

	//private()
}
