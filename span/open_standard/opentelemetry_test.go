package open_standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
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

	rootSpan := field.NewSpanFromPool(nil, "")

	// 1.0 set external traceID and parentID
	rootSpan.SetParentID(field.GenSpanID())
	rootSpan.SetTraceID(field.GenTraceID())

	attrs := field.MallocStructField(3)
	attrs.Set("work", field.StringField("test"))
	attrs.Set("testFunc", field.StringField("TestOpenTelemetryWrite"))
	attrs.Set("testSpan", field.StringField("root"))
	rootSpan.SetAttributes("test", attrs)

	// 1.1 record event to span
	setTestSpance(rootSpan)
	b := bytes.NewBuffer(nil)

	// 1.2 record metrics to span
	rootSpan.Metric(getTestMetric())
	rootSpan.Metric(getTestMetric())
	rootSpan.Metric(getTestMetric())

	// 1.3 sub task or thread need new children Span
	childrenSpan0 := rootSpan.Children()

	// 1.4 Get external span from internalSpan
	external0 := rootSpan.NewExternalSpan()
	// 1.4.1 record info to external span
	external0.StartTime = time.Now()
	external0.EndTime = time.Now()
	// external0.

	// 1.5. won't use root span, free it
	rootSpan.Signal()

	// 2.1 record sub span
	setTestSpance(childrenSpan0)

	// 2.2 write external span
	es0 := childrenSpan0.NewExternalSpan()
	es0.Attributes.Set("test.opentelemetry", field.StringField("telemetry"))
	es0.Attributes.Set("method", field.StringField("test telemetry"))

	// 2.2 won't use sub span, free sub span
	childrenSpan0.Signal()

	enc := encoder.NewJsonEncoder(b)
	open := OpenTelemetry{
		Encoder: enc,
	}
	open.SetDefultResources()
	err := open.Write(rootSpan)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	enc.Close()

	// check result
	cap := map[string]interface{}{}
	bytes := b.Bytes()
	left := 0
	i := 0
	n := 0
	for ; i < len(bytes); i += 1 {
		if bytes[i] == '\n' {
			if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
				t.Error(err)
				t.FailNow()
			} else {
				n += 1
			}
			left = i + 1
		}
	}
	if left < len(bytes) {
		if err = json.Unmarshal(bytes[left:i], &cap); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			n += 1
		}
	}

	fmt.Print(b.String())
}
