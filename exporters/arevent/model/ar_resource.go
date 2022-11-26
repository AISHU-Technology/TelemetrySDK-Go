package model

// ARResource ，记录资源信息例如服务名、版本号、主机信息等。
type ARResource interface {
	// GetSchemaURL 返回 SchemaURL 。
	GetSchemaURL() string
	// GetAttributes 返回 Attributes 。
	GetAttributes() map[string]interface{}
}
