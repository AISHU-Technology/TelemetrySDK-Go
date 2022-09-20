package common

import "go.opentelemetry.io/otel/attribute"

// Attribute 自定义Attribute统一Type。
type Attribute struct {
	Key   string
	Value Value
}

// Value 自定义Value统一Type。
type Value struct {
	Type  string
	Value interface{}
}

// AnyRobotAttributeFromKeyValue 单条KeyValue转换为*Attribute。
func AnyRobotAttributeFromKeyValue(kv attribute.KeyValue) *Attribute {
	return &Attribute{
		Key: string(kv.Key),
		Value: Value{
			Type:  standardizeValueType(kv.Value.Type().String()),
			Value: kv.Value.AsInterface(),
		},
	}
}

// AnyRobotAttributesFromKeyValues 批量KeyValue转换为[]*Attribute。
func AnyRobotAttributesFromKeyValues(kvs []attribute.KeyValue) []*Attribute {
	if kvs == nil {
		return make([]*Attribute, 0, 0)
	}
	attributes := make([]*Attribute, 0, len(kvs))
	for i := 0; i < len(kvs); i++ {
		attributes = append(attributes, AnyRobotAttributeFromKeyValue(kvs[i]))
	}
	return attributes
}

// standardizeValueType 标准化统一 ValueType 为各语言统一格式。
func standardizeValueType(valueType string) string {
	switch valueType {
	case "BOOLSLICE":
		return "BOOLARRAY"
	case "INT64":
		return "INT"
	case "INT64SLICE":
		return "INTARRAY"
	case "FLOAT64":
		return "FLOAT"
	case "FLOAT64SLICE":
		return "FLOATARRAY"
	case "STRINGSLICE":
		return "STRINGARRAY"
	default:
		return valueType
	}
}
