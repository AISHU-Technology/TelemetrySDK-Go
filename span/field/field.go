package field

import "time"

type IntField int
type Float64Field float64
type StringField string
type TimeField time.Time
type FieldTpye int
type JsonFiled struct {
	Data interface{}
}

const (
	IntType = iota
	Float64Type
	StringType
	TimeType
	ArrayType
	StructType
	MetricType
	JsonType
)

type Field interface {
	Type() FieldTpye
	protect()
}

func (f IntField) Type() FieldTpye {
	return IntType
}

func (f IntField) protect() {}

func (f Float64Field) Type() FieldTpye {
	return Float64Type
}

// Avoiding irrelevant personnel to implement Field interface
func (f Float64Field) protect() {}

func (f StringField) Type() FieldTpye {
	return StringType
}

// Avoiding irrelevant personnel to implement Field interface
func (f StringField) protect() {}

func (f TimeField) Type() FieldTpye {
	return TimeType
}

// Avoiding irrelevant personnel to implement Field interface
func (f TimeField) protect() {}

func (f *ArrayField) Type() FieldTpye {
	return ArrayType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *ArrayField) protect() {}

func (f *StructField) Type() FieldTpye {
	return StructType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *StructField) protect() {}

func (f *JsonFiled) Type() FieldTpye {
	return JsonType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *JsonFiled) protect() {}

func MallocJsonField(data interface{}) *JsonFiled {
	return &JsonFiled{
		Data: data,
	}
}
