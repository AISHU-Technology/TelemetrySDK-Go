package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	"testing"
	"time"

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

	// r1 := testLogOne{
	// 	Level: 1,
	// 	// Message: "",
	// }

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.Record(r1)

	m0 := getTestMetric()

	s.Metric(m0)

	// t0 := TraceSpan{}
	// setTestTrace(&t0)
	// s.Trace(&t0)
}

func TestRecord(t *testing.T) {
	// runtimeSpan := NewRuntime(NewSpanFromPool)
	buf := bytes.NewBuffer(nil)
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}
	runtimeSpan := NewRuntime(writer, field.NewSpanFromPool)
	// go runtimeSpan.Run()

	// task thread0, log quick
	go func() {
		// a new web request task
		root := runtimeSpan.Children()

		// web request complete, wait children span
		// finally send web's span control to loger runtime
		defer root.Signal()

		// do some task and log
		setTestSpance(root)

		rl := root.ListRecord()
		assert.Equal(t, 2, len(rl))

		// a new children thread or a sub task start
		sr0 := root.Children()

		// task thread0's sub thread, we call thread2 live 2s after task thread0 complete
		go func() {

			// task complete, send span control to parent
			defer sr0.Signal()

			// do some task and log
			setTestSpance(sr0)

			time.Sleep(2 * time.Second)

			// rl = sr0.ListRecord()
			// assert.Equal(t, 2, len(rl))
		}()

	}()

	// task thread1, live 3s
	go func() {
		root := runtimeSpan.Children()
		defer root.Signal()
		setTestSpance(root)

		time.Sleep(3 * time.Second)

	}()

	// stop runtime after 4s
	go func() {
		time.Sleep(4 * time.Second)
		runtimeSpan.Signal()
	}()

	runtimeSpan.Run()

	// test runtime stop
	after_stop := runtimeSpan.Children()
	assert.Equal(t, nil, after_stop)

	// check result
	var err error
	cap := map[string]interface{}{}
	bytes := buf.Bytes()
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

	assert.Equal(t, 3, n)

}

func TestSignal(t *testing.T) {
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(ioutil.Discard),
	}
	runtimeSpan := NewRuntime(writer, field.NewSpanFromPool)

	runtimeSpan.Signal()
	runtimeSpan.Signal()
}
