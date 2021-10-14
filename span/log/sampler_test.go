package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testLogString(l *SamplerLogger) {
	attr := field.NewAttribute("test", field.StringField("test attr"))
	l.Trace(string(TraceLevelString), nil)
	l.Debug(string(DebugLevelString), nil)
	l.Info(string(InfoLevelString), nil)
	l.Warn(string(WarnLevelString), nil)
	l.Error(string(ErrorLevelString), nil)
	l.Fatal(string(FatalLevelString), nil)

	l.Trace(string(TraceLevelString), attr)
	l.Debug(string(DebugLevelString), attr)
	l.Info(string(InfoLevelString), attr)
	l.Warn(string(WarnLevelString), attr)
	l.Error(string(ErrorLevelString), attr)
	l.Fatal(string(FatalLevelString), attr)
}

func testLogField(l *SamplerLogger) {
	attr := field.NewAttribute("test", field.StringField("test attr"))
	l.TraceField(TraceLevelString, "test", attr)
	l.DebugField(DebugLevelString, "test", attr)
	l.InfoField(InfoLevelString, "test", attr)
	l.WarnField(WarnLevelString, "test", attr)
	l.ErrorField(ErrorLevelString, "test", attr)
	l.FatalField(FatalLevelString, "test", attr)

	l.TraceField(TraceLevelString, "test", nil)
	l.DebugField(DebugLevelString, "test", nil)
	l.InfoField(InfoLevelString, "test", nil)
	l.WarnField(WarnLevelString, "test", nil)
	l.ErrorField(ErrorLevelString, "test", nil)
	l.FatalField(FatalLevelString, "test", nil)
}

func testLogLevel(t *testing.T, l *SamplerLogger, level int) {
	l.SetLevel(level)
	testLogString(l)
	testLogField(l)
}

func TestSamplerLoggerSpan(t *testing.T) {
	buf := ioutil.Discard
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	l.SetContext(context.Background())

	testLogLevel(t, l, TraceLevel)
	testLogLevel(t, l, DebugLevel)
	testLogLevel(t, l, InfoLevel)
	testLogLevel(t, l, WarnLevel)
	testLogLevel(t, l, ErrorLevel)
	testLogLevel(t, l, FatalLevel)

	l.Close()
}

func TestSampleCheck(t *testing.T) {
	count := 1000000
	total := 0
	l := NewDefaultSamplerLogger()
	l.Sample = 1
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}
	assert.Equal(t, count, total)

	l.Sample = 0
	total = 0
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}
	assert.Equal(t, 0, total)

	l.Sample = 0.8
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}

	// allow +- 0.5
	if total > int(float32(count)*(l.Sample+0.5)) || total < int(float32(count)*(l.Sample-0.5)) {
		t.Errorf("TestSampleCheck error, total: %d, sample: %d", total, int(float32(count)*l.Sample))
	}
}

func TestSamplerLoggerNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// test nil span to record
	testLogField(l)

	l.Close()
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() > 0, "record shouldn't drop")

	// test runtime is nil
	run = runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	l.runtime = nil
	buf.Reset()
	testLogField(l)
	l.Close()
	run.Signal()
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() == 0, "record should drop")

	testLogString(l)
	assert.True(t, buf.Len() == 0, "record should drop")

}

func TestSamplerLoggerClose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	enc := encoder.NewJsonEncoder(buf)
	ot := open_standard.NewOpenTelemetry(enc, nil)
	run := runtime.NewRuntime(&ot, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	// runtime will wait sub LogSpan completed

	go l.Close()
	testLogField(l)

	// run.Signal()
	// assert.Equal(t, nil, enc.Close())
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() > 0, "record drop")

	buf.Reset()

	// will not log after runtime is closed
	testLogField(l)

	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() == 0, "record error")
}

func TestSamplerLogger(t *testing.T) {
	// 0. create logger and start runtime
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// 1.1 log message into roor LogSpan
	l.Debug("debug string message", nil)
	l.DebugField(field.StringField("debug field message"), "test", nil)

	attrs := field.MallocStructField(3)
	attrs.Set("work", field.StringField("test"))
	attrs.Set("testFunc", field.StringField("TestSamplerLogger"))
	attrs.Set("testSpan", field.StringField("root"))
	// set attr

	attr := field.NewAttribute("attr", attrs)
	l.Info("infomessage", attr)

	// final close runtime and clean work space
	l.Close()
	// run.Signal()

	time.Sleep(1 * time.Second)

	cap := map[string]interface{}{}
	bytes := buf.Bytes()
	left := 0
	i := 0
	n := 0
	for ; i < len(bytes); i += 1 {
		if bytes[i] == '\n' {
			if err := json.Unmarshal(bytes[left:i], &cap); err != nil {
				t.Error(err)
				t.FailNow()
			} else {
				n += 1
				fmt.Println(string(bytes[left:i]))
				fmt.Println()
			}
			left = i + 1
		}
	}
	if left < len(bytes) {
		if err := json.Unmarshal(bytes[left:i], &cap); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			n += 1
			fmt.Println(string(bytes[left:i]))
		}
	}

	// fmt.Print(buf.String())
}
