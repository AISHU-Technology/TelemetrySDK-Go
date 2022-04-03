package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
	"time"

	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/field"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/open_standard"

	"github.com/stretchr/testify/assert"
)

func fakeTestStructField() field.Field {
	_, msg, line, _ := runtime.Caller(1)

	res := field.MallocStructField(4)
	res.Set("Level", field.IntField(0))
	res.Set("Model", field.StringField("test Eacape \\\"\b\f\\t{}\r\n"))
	res.Set("Message", field.StringField(fmt.Sprintf("%s: %d", msg, line)))
	res.Set("eventNum", field.IntField(1))

	return res
}

func setTestSpance(s field.LogSpan) {
	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.SetRecord(r1)

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
		root := runtimeSpan.Children(nil)

		// web request complete, wait children span
		// finally send web's span control to loger runtime
		defer root.Signal()

		// do some task and log
		setTestSpance(root)
	}()
	// task thread1, live 1s
	go func() {
		root := runtimeSpan.Children(nil)
		defer root.Signal()
		setTestSpance(root)

		time.Sleep(1 * time.Second)

	}()

	// stop runtime after 2s
	go func() {
		time.Sleep(2 * time.Second)
		runtimeSpan.Signal()
	}()

	runtimeSpan.Run()

	// test runtime stop
	after_stop := runtimeSpan.Children(nil)
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

	assert.Equal(t, 2, n)

}

func TestSignal(t *testing.T) {
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(ioutil.Discard),
	}
	runtimeSpan := NewRuntime(writer, field.NewSpanFromPool)

	runtimeSpan.Signal()
	runtimeSpan.Signal()
}
