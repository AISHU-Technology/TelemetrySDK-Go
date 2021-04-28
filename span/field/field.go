package field

import "time"

type IntField int
type Float64Field float64
type StringField string
type TimeField time.Time
type FieldTpye int

const (
	IntType = iota
	Float64Type
	StringType
	TimeType
	ArrayType
	StructType
	ExternalSpanType
	MetricType
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

func (f Float64Field) protect() {}

func (f StringField) Type() FieldTpye {
	return StringType
}

func (f StringField) protect() {}

func (f TimeField) Type() FieldTpye {
	return TimeType
}

func (f TimeField) protect() {}

func (f *ArrayField) Type() FieldTpye {
	return ArrayType
}

func (f *ArrayField) protect() {}

func (f *StructField) Type() FieldTpye {
	return StructType
}

func (f *StructField) protect() {}

func (f *Mmetric) Type() FieldTpye {
	return MetricType
}

func (f *Mmetric) protect() {}

func (f *ExternalSpanField) Type() FieldTpye {
	return ExternalSpanType
}

func (f *ExternalSpanField) protect() {}
