package common

import "go.opentelemetry.io/otel/attribute"

// Attribute 自定义Attribute统一Type。
type Attribute struct {
	Key   string `json:"Key"`
	Value Value  `json:"Value"`
}

// Value 自定义Value统一Type。
type Value struct {
	Type  string      `json:"Type"`
	Value interface{} `json:"Value"`
}

// AnyRobotAttributeFromKeyValue 单条KeyValue转换为*Attribute。
func AnyRobotAttributeFromKeyValue(attribute attribute.KeyValue) *Attribute {
	return &Attribute{
		Key: string(attribute.Key),
		Value: Value{
			Type:  standardizeValueType(attribute.Value.Type().String()),
			Value: attribute.Value.AsInterface(),
		},
	}
}

// AnyRobotAttributesFromKeyValues 批量KeyValue转换为[]*Attribute。
func AnyRobotAttributesFromKeyValues(attributes []attribute.KeyValue) []*Attribute {
	if attributes == nil {
		return make([]*Attribute, 0)
	}
	arattributes := make([]*Attribute, 0, len(attributes))
	for i := 0; i < len(attributes); i++ {
		arattributes = append(arattributes, AnyRobotAttributeFromKeyValue(attributes[i]))
	}
	return arattributes
}

// standardizeValueType 标准化统一 ValueType 为各语言统一格式。
func standardizeValueType(valueType string) string {
	switch valueType {
	//case "BOOL":
	//	return "BOOL"
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
	//case "STRING":
	//	return "STRING"
	case "STRINGSLICE":
		return "STRINGARRAY"
	default:
		return valueType
	}
}
