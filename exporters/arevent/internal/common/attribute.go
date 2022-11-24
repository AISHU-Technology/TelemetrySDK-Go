package common

import "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"

// Attribute 自定义 Event Attribute 和 Trace 输出格式一致。
type Attribute struct {
	Key   string        `json:"Key"`
	Value model.ARValue `json:"Value"`
}

// Value 自定义Value统一Type为8种枚举类型。
type Value struct {
	Type  string      `json:"Type"`
	Value interface{} `json:"Value"`
}

// NewAttribute 创建新的 Attribute 。
func NewAttribute(key string, value model.ARValue) model.ARAttribute {
	return &Attribute{
		Key:   key,
		Value: value,
	}
}

func (a *Attribute) Valid() bool {
	return a.keyDefined() && a.valueTyped()
}

// keyDefined 校验 Attribute.Key 不为空，即有含义。
func (a *Attribute) keyDefined() bool {
	return len(a.Key) > 0
}

// valueTyped 校验 Value.Type 是枚举类型。
func (a *Attribute) valueTyped() bool {
	switch a.Value.GetType() {
	case "BOOL":
		return true
	case "BOOLARRAY":
		return true
	case "INT":
		return true
	case "INTARRAY":
		return true
	case "FLOAT":
		return true
	case "FLOATARRAY":
		return true
	case "STRING":
		return true
	case "STRINGARRAY":
		return true
	default:
		return false
	}
}

func (a *Attribute) GetKey() string {
	return a.Key
}

func (a *Attribute) GetValue() model.ARValue {
	return a.Value
}

func (v Value) GetType() string {
	return v.Type
}

func (v Value) GetValue() interface{} {
	return v.Value
}

// BoolValue 传入 bool 类型的值。
func BoolValue(value bool) model.ARValue {
	return Value{
		Type:  "BOOL",
		Value: value,
	}
}

// BoolArray 传入 []bool 类型的值。
func BoolArray(value []bool) model.ARValue {
	return Value{
		Type:  "BOOLARRAY",
		Value: value,
	}
}

// IntValue 传入 int 类型的值。
func IntValue(value int) model.ARValue {
	return Value{
		Type:  "INT",
		Value: value,
	}
}

// IntArray 传入 []int 类型的值。
func IntArray(value []int) model.ARValue {
	return Value{
		Type:  "INTARRAY",
		Value: value,
	}
}

// FloatValue 传入 float64 类型的值。
func FloatValue(value float64) model.ARValue {
	return Value{
		Type:  "FLOAT",
		Value: value,
	}
}

// FloatArray 传入 []float64 类型的值。
func FloatArray(value []float64) model.ARValue {
	return Value{
		Type:  "FLOATARRAY",
		Value: value,
	}
}

// StringValue 传入 string 类型的值。
func StringValue(value string) model.ARValue {
	return Value{
		Type:  "STRING",
		Value: value,
	}
}

// StringArray 传入 []string 类型的值。
func StringArray(value []string) model.ARValue {
	return Value{
		Type:  "STRINGARRAY",
		Value: value,
	}
}
