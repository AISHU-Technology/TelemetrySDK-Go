package field

import (
	"strconv"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestArrayField(t *testing.T) {
	cap := 10
	length := 11
	f := MallocArrayField(cap)

	for i := 0; i < length; i += 1 {
		f.Append(IntField(i))
	}

	for i := 0; i < length; i += 1 {
		v := IntField(i)
		assert.Equal(t, v, (*f)[i])
		v1, err := f.At(i)
		assert.Equal(t, v, v1)
		assert.Equal(t, err, nil)
	}

	fs := MallocStructField(1)
	f.Append(fs)
	assert.Equal(t, (*f)[length], fs)

	fa := MallocArrayField(0)
	f.Append(fa)
	assert.Equal(t, (*f)[length+1], fa)

	assert.Equal(t, f.Length(), length+2)

	v, err := f.At(length + 3)
	assert.Equal(t, err, OverIndexError)
	assert.Equal(t, v, nil)

}

func TestStructField(t *testing.T) {
	cap := 10
	length := 11
	f := MallocStructField(cap)
	for i := 0; i < length; i += 1 {
		f.Set(strconv.Itoa(i), IntField(i))
	}

	for i := 0; i < length; i += 1 {
		k := strconv.Itoa(i)
		v := IntField(i)
		assert.Equal(t, v, f.values[i])
		assert.Equal(t, k, f.keys[i])
		k1, v1, err := f.At(i)
		assert.Equal(t, k1, k)
		assert.Equal(t, v1, v)
		assert.Equal(t, err, nil)
	}

	fs := MallocStructField(1)
	fs.Set("test", StringField("sting"))
	f.Set("struct", fs)
	assert.Equal(t, f.keys[length], "struct")
	assert.Equal(t, f.values[length], fs)

	fa := MallocArrayField(0)
	f.Set("array", fa)
	assert.Equal(t, f.keys[length+1], "array")
	assert.Equal(t, f.values[length+1], fa)

	assert.Equal(t, f.Length(), length+2)

	k, v, err := f.At(length + 3)
	assert.Equal(t, k, "")
	assert.Equal(t, v, nil)
	assert.Equal(t, err, OverIndexError)
}

func TestFieldType(t *testing.T) {
	assert.Equal(t, FieldTpye(IntType), IntField(0).Type())
	assert.Equal(t, FieldTpye(Float64Type), Float64Field(0).Type())
	assert.Equal(t, FieldTpye(StringType), StringField("").Type())
	assert.Equal(t, FieldTpye(TimeType), TimeField(time.Now()).Type())
	assert.Equal(t, FieldTpye(ArrayType), MallocArrayField(0).Type())
	assert.Equal(t, FieldTpye(StructType), MallocStructField(0).Type())
	assert.Equal(t, FieldTpye(ExternalSpanType), (&ExternalSpanField{}).Type())
	assert.Equal(t, FieldTpye(MetricType), (&Mmetric{}).Type())
	assert.Equal(t, FieldTpye(JsonType), (&JsonFiled{}).Type())
}

