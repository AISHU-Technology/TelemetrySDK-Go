package open_standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"

	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"

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
	r0 := fakeTestStructField()
	s.SetRecord(r0)

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.SetRecord(r1)
	record := field.MallocStructField(2)
	record.Set("test", r1)
	//attr := &field.Attribute{
	//	Type:    "test",
	//	Message: record,
	//}
	//
	//s.SetAttributes(attr)

}

func TestOpenTelemetryWrite(t *testing.T) {

	rootSpan := field.NewSpanFromPool(nil, nil)

	setTestSpance(rootSpan)
	b := bytes.NewBuffer(nil)

	// 1.5. won't use root span, free it
	rootSpan.Signal()

	enc := encoder.NewJsonEncoder(b)
	open := OpenTelemetry{
		Encoder:  enc,
		Resource: nil,
	}

	defer open.Close()
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

func TestOpenTelemetrySetDefaultResources(t *testing.T) {
	b := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(b)
	open := NewOpenTelemetry(enc, nil)
	os.Setenv("HOSTNAME", "test")
	f := field.MallocStructField(10)
	f.Set("HOSTNAME", field.StringField("test"))
	f.Set("Telemetry.SDK.Name", field.StringField(SDKName))
	f.Set("Telemetry.SDK.Version", field.StringField(SDKVersion))
	f.Set("Telemetry.SDK.Language", field.StringField(SDKLanguage))

	open.SetDefaultResources()
	assert.Equal(t, open.Resource, f)
}
