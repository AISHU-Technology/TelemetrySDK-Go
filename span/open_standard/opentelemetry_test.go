package open_standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"span/encoder"
	"span/field"
	"testing"
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

// func setTestTrace(t *TraceSpan) {
// 	t.Header = http.Header{
// 		"model": []string{"test"},
// 	}
// 	t.Method = "GET"
// 	t.Client = endpoint{
// 		Port: "0",
// 		Addr: "127.0.0.1",
// 	}
// 	t.Uri = "/test"
// 	t.StatusCode = 200
// 	t.StartTime = time.Now()
// 	t.EndTime = time.Now().Add(120)
// 	t.Type = "http"
// }

func fakeTestStructField() field.Field {
	_, msg, line, _ := runtime.Caller(1)

	res := field.MallocStructField(4)
	res.Set("Level", field.IntField(0))
	res.Set("Model", field.StringField("test Eacape \\\"\b\f\\t{}\r\n"))
	res.Set("Message", field.StringField(fmt.Sprintf("%s: %d", msg, line)))
	res.Set("eventNum", field.IntField(1))

	return res
}

func setTestSpance(s field.InternalSpan) {
	r0 := fakeTestStructField()
	s.Record(r0)

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.Record(r1)

	m0 := getTestMetric()

	s.Metric(m0)
}

func TestOpenTelemetryWrite(t *testing.T) {
	s := field.NewSpanFromPool(nil, "")
	setTestSpance(s)
	b := bytes.NewBuffer(nil)

	s.Metric(getTestMetric())
	s.Metric(getTestMetric())
	s.Metric(getTestMetric())

	c0 := s.Children()
	setTestSpance(c0)
	c0.Signal()

	enc := encoder.NewJsonEncoder(b)
	open := OpenTelemetry{
		Encoder: enc,
	}
	err := open.Write(s)
	if err != nil {
		t.Error(err)
	}

	// check result
	cap := map[string]interface{}{}
	bytes := b.Bytes()
	left := 0
	i := 0
	for ; i < len(bytes); i += 1 {
		if bytes[i] == '\n' {
			if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
				t.Error(err)
			}
			left = i + 1
		}
	}
	if left < len(bytes) {
		if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
			t.Error(err)
		}
	}

	fmt.Println("--------log-------------:")
	fmt.Print(b.String())
}
