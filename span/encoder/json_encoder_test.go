package encoder

import (
	bytess "bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"testing"

	"gotest.tools/assert"
)

func getTestMetric() field.Mmetric {
	m := field.Mmetric{}
	m.AddLabel("test")
	m.AddLabel("0")
	m.AddLabel("Metric")
	m.Set("testMetric0", 0.0)
	m.AddAttribute("testAttr", "test")
	return m
}

func GetTestFieds() []field.Field {
	r0 := fakeTestStructField()

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))

	return []field.Field{r0, r1}

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

	// test metric encoder
	b.Reset()
	metrics := []field.Mmetric{}
	metrics = append(metrics, getTestMetric())
	metrics = append(metrics, getTestMetric())
	metrics = append(metrics, getTestMetric())

	for _, i := range metrics {
		enc.Write(&i)
	}
	enc.Close()

	// fmt.Println("--------metric-------------:")
	left = 0
	i = 0
	bytess = b.Bytes()
	for ; i < len(bytess); i += 1 {
		if bytess[i] == '\n' {
			if err = json.Unmarshal(bytess[left:i], &check); err != nil {
				t.Error(err)
			}
			left = i + 1
		}
	}
	fmt.Println(b.String())

	// test external Span encoder
	b.Reset()
	s := field.NewSpanFromPool(nil, "")
	s.SetTraceID(field.GenTraceID())
	s.SetParentID(field.GenSpanID())
	es0 := s.NewExternalSpan()
	es0.Attributes.Set("method", field.StringField("test"))
	es1 := s.NewExternalSpan()
	es1.Attributes.Set("host", field.StringField("test"))
	s.Signal()

	for _, i := range s.ListExternalSpan() {
		enc.Write(i)
	}

	// fmt.Println("--------extrenalSpan-------------:")
	left = 0
	i = 0
	bytess = b.Bytes()
	for ; i < len(bytess); i += 1 {
		if bytess[i] == '\n' {
			// fmt.Println(string(bytess[left:i]))
			if err = json.Unmarshal(bytess[left:i], &check); err != nil {
				t.Error(err)
			}
			left = i + 1
		}
	}
	fmt.Println(b.String())
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
