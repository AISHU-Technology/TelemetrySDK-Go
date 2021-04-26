package encoder

import (
	"bytes"
	"fmt"
	"runtime"
	"span/field"
	"testing"
	"time"
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
	res.Set("Model", field.StringField("test Eacape \\\"\b\f\\t{}\r\n"))
	res.Set("Message", field.StringField(fmt.Sprintf("%s: %d", msg, line)))
	res.Set("eventNum", field.IntField(1))

	return res
}

func TestJsonEncoder(t *testing.T) {
	fields := GetTestFieds()
	b := bytes.NewBuffer(nil)
	nowTime := time.Now()
	fields = append(fields, field.TimeField(nowTime))

	enc := NewJsonEncoder(b)
	for _, i := range fields {
		if err := enc.Write(i); err != nil {
			t.Error(err)
		}
	}

	fmt.Println("--------log-------------:")
	fmt.Print(b.String())

	b.Reset()
	metrics := []field.Mmetric{}
	metrics = append(metrics, getTestMetric())
	metrics = append(metrics, getTestMetric())
	metrics = append(metrics, getTestMetric())

	for _, i := range metrics {
		enc.Write(&i)
	}

	fmt.Println("--------metric-------------:")
	fmt.Println(b.String())
}
