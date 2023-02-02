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
type MapField map[string]interface{}

const (
	IntType = iota
	Float64Type
	StringType
	TimeType
	ArrayType
	StructType
	MetricType
	JsonType
	MapType
)

type Field interface {
	Type() FieldTpye
	protect()
}

func (f IntField) Type() FieldTpye {
	return IntType
}

func (f IntField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f Float64Field) Type() FieldTpye {
	return Float64Type
}

// Avoiding irrelevant personnel to implement Field interface
func (f Float64Field) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f StringField) Type() FieldTpye {
	return StringType
}

// Avoiding irrelevant personnel to implement Field interface
func (f StringField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f TimeField) Type() FieldTpye {
	return TimeType
}

// Avoiding irrelevant personnel to implement Field interface
func (f TimeField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *ArrayField) Type() FieldTpye {
	return ArrayType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *ArrayField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *StructField) Type() FieldTpye {
	return StructType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *StructField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *JsonFiled) Type() FieldTpye {
	return JsonType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *JsonFiled) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func MallocJsonField(data interface{}) *JsonFiled {
	return &JsonFiled{
		Data: data,
	}
}

func (f MapField) Type() FieldTpye {
	return MapType
}

func (f MapField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f MapField) Append(key string, value Field) {
	if f == nil {
		f = MallocMapField()
	}
	f[key] = value
}

func MallocMapField() MapField {
	return MapField(make(map[string]interface{}))
}

func WithServiceInfo(ServiceName string, ServiceVersion string, ServiceInstanceID string) Field {
	service := make(map[string]interface{})
	service["name"] = ServiceName
	service["version"] = ServiceVersion
	service["instance"] = map[string]string{"id": ServiceInstanceID}

	mapServiceInfo := MapField(map[string]interface{}{"service": service})
	return mapServiceInfo
}
