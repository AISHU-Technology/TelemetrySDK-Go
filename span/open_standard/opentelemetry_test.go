package open_standard

import (
	"bytes"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	"encoding/json"
	"fmt"
	"testing"
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
	buf := bytes.NewBuffer(nil)
	open := OpenTelemetryWriter(
		encoder.NewJsonEncoder(buf),
		field.IntField(0))

	defer open.Close()

	err := open.Write(rootSpans)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// check result
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
		} // else {
		//n += 1
		//}
	}

	fmt.Print(buf.String())
}
