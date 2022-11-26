package common

// Value 自定义Value统一Type为8种枚举类型。
type Value struct {
	Type string      `json:"Type"`
	Data interface{} `json:"Data"`
}

func (v Value) GetType() string {
	return v.Type
}

func (v Value) GetData() interface{} {
	return v.Data
}

// BoolValue 传入 bool 类型的值。
func BoolValue(value bool) Value {
	return Value{
		Type: "BOOL",
		Data: value,
	}
}

// BoolArray 传入 []bool 类型的值。
func BoolArray(value []bool) Value {
	return Value{
		Type: "BOOLARRAY",
		Data: value,
	}
}

// IntValue 传入 int 类型的值。
func IntValue(value int) Value {
	return Value{
		Type: "INT",
		Data: value,
	}
}

// IntArray 传入 []int 类型的值。
func IntArray(value []int) Value {
	return Value{
		Type: "INTARRAY",
		Data: value,
	}
}

// FloatValue 传入 float64 类型的值。
func FloatValue(value float64) Value {
	return Value{
		Type: "FLOAT",
		Data: value,
	}
}

// FloatArray 传入 []float64 类型的值。
func FloatArray(value []float64) Value {
	return Value{
		Type: "FLOATARRAY",
		Data: value,
	}
}

// StringValue 传入 string 类型的值。
func StringValue(value string) Value {
	return Value{
		Type: "STRING",
		Data: value,
	}
}

// StringArray 传入 []string 类型的值。
func StringArray(value []string) Value {
	return Value{
		Type: "STRINGARRAY",
		Data: value,
	}
}
