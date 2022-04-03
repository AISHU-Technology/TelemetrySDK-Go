package encoder

import (
	bytess "bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/field"

	"github.com/stretchr/testify/assert"
)

func getTestJson() field.Field {

	type A struct {
		Name string
		Age  int
	}
	var a = &A{Name: "123", Age: 12}

	return field.MallocJsonField(a)
}

func GetTestFieds() []field.Field {
	r0 := fakeTestStructField()

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))

	r3 := getTestJson()
	return []field.Field{r0, r1, r3}

	// m0 := getTestMetric()

	// s.Metric(m0)

	// t0 := TraceSpan{}
	// setTestTrace(&t0)
	// s.Trace(&t0)
}

func fakeTestStructField() field.Field {
	_, msg, line, _ := runtime.Caller(1)

	res := field.MallocStructField(4)
	res.Set("Level", field.IntField(0))
	res.Set("Model", field.StringField("test Eacape \\\"\b\f\\t\t{}\r\n"))
	res.Set("Message", field.StringField(fmt.Sprintf("%s: %d", msg, line)))
	res.Set("eventNum", field.IntField(1))
	res.Set("float", field.Float64Field(1.3))

	return res
}

func TestJsonEncoder(t *testing.T) {
	check := map[string]interface{}{}
	var err error

	// test log encoder
	fields := GetTestFieds()

	b := bytess.NewBuffer(nil)
	enc := NewJsonEncoder(b)
	for _, i := range fields {
		if err := enc.Write(i); err != nil {
			t.Error(err)
		}
	}

	// fmt.Println("--------log-------------:")
	left := 0
	i := 0
	bytess := b.Bytes()
	for ; i < len(bytess); i += 1 {
		if bytess[i] == '\n' {
			if err = json.Unmarshal(bytess[left:i], &check); err != nil {
				t.Error(err)
			}
			left = i + 1
		}
	}
	fmt.Print(b.String())

}

func TestArrayField(t *testing.T) {
	cap := 10
	length := 11
	a := field.MallocArrayField(cap)

	for i := 0; i < length; i += 1 {
		a.Append(field.IntField(i))
	}

	for i := 0; i < length; i += 1 {
		assert.Equal(t, field.IntField(i), (*a)[i])
	}

	b := bytess.NewBuffer(nil)
	enc := NewJsonEncoder(b)
	enc.Write(a)

}

func TestNewJsonEncoderBench(t *testing.T) {
	b := bytess.NewBuffer(nil)
	en := NewJsonEncoderBench(b)
	err := en.Write(field.TimeField(time.Now()))
	if err != nil {
		panic(err)
	}
	en.Close()

}
