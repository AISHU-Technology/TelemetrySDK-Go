package model

type ARResource interface {
	GetSchemaURL() string
	GetAttributes() []*ARAttribute
	//private()
}
