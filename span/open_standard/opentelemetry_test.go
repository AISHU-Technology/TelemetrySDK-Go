package open_standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"

	"github.com/stretchr/testify/assert"
)

func setTestSpance(s field.LogSpan) {
	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.SetRecord(r1)
}

func TestOpenTelemetryWrite(t *testing.T) {

	rootSpan := field.NewSpanFromPool(nil, nil)

	setTestSpance(rootSpan)
	var rootSpans []field.LogSpan
	rootSpans = append(rootSpans, rootSpan)
	b := bytes.NewBuffer(nil)

	enc := encoder.NewJsonEncoder(b)
	open := OpenTelemetry{
		Encoder:  enc,
		Resource: nil,
	}

	defer open.Close()

	err := open.Write(rootSpans)
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
		} // else {
		//n += 1
		//}
	}

	fmt.Print(b.String())
}

func TestOpenTelemetrySetDefaultResources(t *testing.T) {
	b := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(b)
	open := NewOpenTelemetry(enc, nil)
	defaultResource := make(map[string]interface{})
	service := make(map[string]interface{})
	service["name"] = serviceName
	service["version"] = serviceVersion
	service["instance"] = map[string]string{"id": serviceInstance}
	defaultResource["service"] = service
	// o.Resource = field.MapField(defaultResource)
	// os.Setenv("HOSTNAME", "test")
	// f := field.MallocStructField(10)
	// f.Set("HOSTNAME", field.StringField("test"))
	// f.Set("Telemetry.SDK.Name", field.StringField(sdkName))
	// f.Set("Telemetry.SDK.Version", field.StringField(sdkVersion))
	// f.Set("Telemetry.SDK.Language", field.StringField(sdkLanguage))

	open.SetDefaultResources()
	assert.Equal(t, open.Resource, field.MapField(defaultResource))
}

func TestOpenTelemetrySetDefaultResourcesWithServiceInfo(t *testing.T) {
	b := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(b)
	open := NewOpenTelemetry(enc, nil)
	defaultResource := make(map[string]interface{})
	service := make(map[string]interface{})
	service["name"] = "testServiceName"
	service["version"] = "testServiceVersion"
	service["instance"] = map[string]string{"id": "testServiceInstanceID"}
	defaultResource["service"] = service
	open.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	assert.Equal(t, open.Resource, field.MapField(defaultResource))
}

func TestDealResource(t *testing.T) {
	b := bytes.NewBuffer(nil)
	enc := encoder.NewJsonEncoder(b)
	open := NewOpenTelemetry(enc, field.IntField(1))
	open.dealResource()
	defaultResource := getDefaultResource()
	service := make(map[string]interface{})
	service["name"] = serviceName
	service["version"] = serviceVersion
	service["instance"] = map[string]string{"id": serviceInstance}
	defaultResource["service"] = service
	defaultResource["customer"] = field.IntField(1)
	assert.Equal(t, open.Resource, field.MapField(defaultResource))
}
