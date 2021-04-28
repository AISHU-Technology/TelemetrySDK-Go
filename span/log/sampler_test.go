package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"span/encoder"
	"span/field"
	"span/open_standard"
	"span/runtime"
	"testing"
	"time"

	"gotest.tools/assert"
)

func testLogString(l *SamplerLogger, span field.InternalSpan) {
	l.Trace(string(TraceLevelString), span)
	l.Debug(string(DebugLevelString), span)
	l.Info(string(InfoLevelString), span)
	l.Warn(string(WarnLevelString), span)
	l.Error(string(ErrorLevelString), span)
	l.Fatal(string(FatalLevelString), span)
}

func testLogField(l *SamplerLogger, span field.InternalSpan) {
	l.TraceField(TraceLevelString, span)
	l.DebugField(DebugLevelString, span)
	l.InfoField(InfoLevelString, span)
	l.WarnField(WarnLevelString, span)
	l.ErrorField(ErrorLevelString, span)
	l.FatalField(FatalLevelString, span)
}

func testLogLevel(t *testing.T, l *SamplerLogger, level int) {
	l.LogLevel = level
	s := l.NewInternalSpan()
	testLogString(l, s)
	s.Signal()
	assert.Equal(t, len(s.ListRecord()), FatalLevel-level+1)

	s = l.NewInternalSpan()
	testLogField(l, s)
	assert.Equal(t, len(s.ListRecord()), FatalLevel-level+1)
	s.Signal()

}

func TestSamplerLoggerSpan(t *testing.T) {
	buf := ioutil.Discard
	l := NewdefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	testLogLevel(t, l, TraceLevel)
	testLogLevel(t, l, DebugLevel)
	testLogLevel(t, l, InfoLevel)
	testLogLevel(t, l, WarnLevel)
	testLogLevel(t, l, ErrorLevel)
	testLogLevel(t, l, FatalLevel)

	l.Close()
}

// func TestDefer(t *testing.T) {
// 	fmt.Println("1111")
// 	b := true
// 	if b {
// 		defer fmt.Println("3333")
// 	}
// 	fmt.Println("2222")
// }

func TestSampleCheck(t *testing.T) {
	count := 1000000
	total := 0
	l := NewdefaultSamplerLogger()
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
	assert.Assert(t, total <= int(float32(count)*l.Sample), "TestSampleCheck error")
}

func TestSamplerLoggerNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewdefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// test nil span to record
	testLogField(l, nil)

	m := field.Mmetric{}
	l.RecordMetrics(m, nil)

	l.Close()
	assert.Assert(t, buf.Len() > 0, "record shouldn't drop")

	// test runtime is nil
	l.runtime = nil
	buf.Reset()
	testLogField(l, nil)
	assert.Assert(t, buf.Len() == 0, "record should drop")

	testLogString(l, nil)
	assert.Assert(t, buf.Len() == 0, "record should drop")

	l.RecordMetrics(m, nil)
	assert.Assert(t, buf.Len() == 0, "metrics should drop")

}

func TestSamplerLoggerClose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewdefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	// runtime will wait sub internalSpan completed
	s := l.NewInternalSpan()
	go l.Close()
	testLogField(l, s)
	s.Signal()

	run.Signal()
	assert.Assert(t, buf.Len() > 0, "record drop1")

	buf.Reset()

	// will not log after runtime is closed
	testLogField(l, nil)
	assert.Assert(t, buf.Len() == 0, "record error")

	assert.Equal(t, l.NewInternalSpan(), nil)
	es, err := l.NewExternalSpan(nil)
	assert.Assert(t, es == nil, "NewExternalSpan error")
	assert.Equal(t, err, field.NilPointerError)

	es, err = l.NewExternalSpan(s)
	assert.Assert(t, es != nil, "NewExternalSpan error")
	assert.Equal(t, err, nil)

}

// func TestwaitgroupSync(t *testing.T) {
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	go func() {
// 		time.Sleep(1 * time.Second)
// 		wg.Done()
// 	}()
// 	wg.Wait()
// 	fmt.Println("wait first")
// 	wg.Wait()
// }
func TestSamplerLogger(t *testing.T) {
	// 0. create logger and start runtime
	buf := bytes.NewBuffer(nil)
	l := NewdefaultSamplerLogger()
	run := runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// 1. first create a root internalSpan
	root := l.NewInternalSpan()

	// 1.0 set trace info for root internalSpan
	traceID := field.GenSpanID()
	externalParentID := field.GenSpanID()
	l.SetTraceID(traceID, root)
	l.SetParentID(externalParentID, root)

	// 1.1 log message into roor internalSpan
	l.Debug("debug string message", root)
	l.DebugField(field.StringField("debug field message"), root)

	// 1.2 create a child internalSpan from root for a sub thread/task
	child0 := l.ChildrenInternalSpan(root)

	// 1.3 start a new thread for sub task
	go func() {
		// 2.1 log message into child internalSpan for child thread
		l.Debug("debug string", child0)

		// 2.X signal child0
		child0.Signal()
	}()

	// 1.4 record some metric into root internalSpan
	m := field.Mmetric{}
	m.Set("root thread", 0.0)
	m.AddLabel("root")
	m.AddLabel("metric")
	m.AddAttribute("root", "root span")
	l.RecordMetrics(m, root)

	// 1.5 record first external request into root internalSpan
	es, err := l.NewExternalSpan(root)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 1.5.1 get trace info for some work
	tID := es.TraceID()
	espID := es.ParentID()
	parentID := es.ParentID()
	spanID := es.ID()
	// 1.5.2 write info to external span
	es.StartTime = time.Now()
	es.EndTime = time.Now()
	es.Attributes.Set("method", field.StringField("test"))
	es.Attributes.Set("host", field.StringField("test"))
	es.Attributes.Set("attr0", field.StringField(tID))
	es.Attributes.Set("attr1", field.StringField(espID))
	es.Attributes.Set("attr2", field.StringField(parentID))
	es.Attributes.Set("attr3", field.StringField(spanID))

	// 1.X signal root internalSpan
	root.Signal()

	// final close runtime and clean work space
	l.Close()
	// run.Signal()

	// check test result
	assert.Equal(t, traceID, tID)
	assert.Equal(t, externalParentID, espID)

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
				fmt.Println(string(bytes[left:i]))
				fmt.Println()
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
			fmt.Println(string(bytes[left:i]))
		}
	}

	// fmt.Print(buf.String())
}
