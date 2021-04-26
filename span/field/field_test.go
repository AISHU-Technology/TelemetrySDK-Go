package field

import (
	"strconv"
	"testing"

	"gotest.tools/assert"
)

func TestArrayField(t *testing.T) {
	cap := 10
	length := 11
	a := MallocArrayField(cap)

	for i := 0; i < length; i += 1 {
		a.Append(IntField(i))
	}

	for i := 0; i < length; i += 1 {
		assert.Equal(t, IntField(i), (*a)[i])
	}

}

func TestStructField(t *testing.T) {
	cap := 10
	length := 11
	f := MallocStructField(cap)
	for i := 0; i < length; i += 1 {
		f.Set(strconv.Itoa(i), IntField(i))
	}

	for i := 0; i < length; i += 1 {
		assert.Equal(t, IntField(i), f.values[i])
		assert.Equal(t, strconv.Itoa(i), f.keys[i])
	}

	f1 := MallocStructField(1)
	f1.Set("test", StringField("sting"))
	f.Set("struct", f1)

	assert.Equal(t, f.keys[length], "struct")
	assert.Equal(t, f.values[length], f1)
}
