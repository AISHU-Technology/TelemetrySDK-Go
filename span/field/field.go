/*
 * @Author: Nick.nie Nick.nie@aishu.cn
 * @Date: 2022-12-09 04:43:20
 * @LastEditors: Nick.nie Nick.nie@aishu.cn
 * @LastEditTime: 2022-12-12 03:15:57
 * @FilePath: /span/field/field.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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

func (f MapField) Type() FieldTpye {
	return MapType
}

func (f MapField) protect() {}

func WithServiceInfo(ServiceName string, ServiceVersion string, ServiceInstanceID string) Field {
	service := make(map[string]interface{})
	service["name"] = ServiceName
	service["version"] = ServiceVersion
	service["instance"] = map[string]string{"id": ServiceInstanceID}

	mapServiceInfo := MapField(map[string]interface{}{"service": service})
	return mapServiceInfo
}
