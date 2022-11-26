package common

// value 自定义Value统一Type为8种枚举类型。
type value struct {
	Type string      `json:"Type"`
	Data interface{} `json:"Data"`
}

func (v value) GetType() string {
	return v.Type
}

func (v value) GetData() interface{} {
	return v.Data
}

func (v value) private() {}

// BoolValue 传入 bool 类型的值。
func BoolValue(v bool) ARValue {
	return value{
		Type: "BOOL",
		Data: v,
	}
}

// BoolArray 传入 []bool 类型的值。
func BoolArray(v []bool) ARValue {
	return value{
		Type: "BOOLARRAY",
		Data: v,
	}
}

// IntValue 传入 int 类型的值。
func IntValue(v int) ARValue {
	return value{
		Type: "INT",
		Data: v,
	}
}

// IntArray 传入 []int 类型的值。
func IntArray(v []int) ARValue {
	return value{
		Type: "INTARRAY",
		Data: v,
	}
}

// FloatValue 传入 float64 类型的值。
func FloatValue(v float64) ARValue {
	return value{
		Type: "FLOAT",
		Data: v,
	}
}

// FloatArray 传入 []float64 类型的值。
func FloatArray(v []float64) ARValue {
	return value{
		Type: "FLOATARRAY",
		Data: v,
	}
}

// StringValue 传入 string 类型的值。
func StringValue(v string) ARValue {
	return value{
		Type: "STRING",
		Data: v,
	}
}

// StringArray 传入 []string 类型的值。
func StringArray(v []string) ARValue {
	return value{
		Type: "STRINGARRAY",
		Data: v,
	}
}
