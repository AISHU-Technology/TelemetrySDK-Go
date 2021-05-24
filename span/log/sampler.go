// Provide some interfaces of logger
//
// looger can use to record string log of any object Log in the form of field
//
// logger can use to log internal thread/task with InteranlSpan.
//
// A InternalSpan use to describe a thread job info or a task info. You can get a thread context info ever recorded.
// And then aggregation Span by TraceID or describe a Trace tree by SpanID and TraceID
package log

import (
	"math/rand"
	"span/field"
	"span/runtime"
	"time"
)

// SamplerLogger implement the Logger interface
// SamplerLogger provide log filter by sampling or log level
//
type SamplerLogger struct {
	// logger sample
	Sample float32
	// logger info
	LogLevel int
	runtime  *runtime.Runtime
}

func NewdefaultSamplerLogger() *SamplerLogger {
	return &SamplerLogger{
		Sample:   1.0,
		LogLevel: InfoLevel,
	}
}

func newStructRecord() *field.StructField {
	return field.MallocStructField(1)
}

// close logger
func (s *SamplerLogger) Close() {
	if s.runtime != nil {
		s.runtime.Signal()
	}
}

// SetRuntime for logger
// this function will signal older runtime, but will not start new runtime
func (s *SamplerLogger) SetRuntime(r *runtime.Runtime) {
	if s.runtime != nil {
		s.runtime.Signal()
	}
	s.runtime = r
}

func (s *SamplerLogger) getInternalSpan() field.InternalSpan {
	if s.runtime != nil {
		return s.runtime.Children()
	}
	return nil
}

func (s *SamplerLogger) sampleCheck() bool {
	if s.Sample >= 1.0 {
		return true
	}

	if s.Sample <= 0 {
		return false
	}

	rand.Seed(time.Now().UnixNano())

	return rand.Float32() <= s.Sample
}

// SetSample set the log sample
// sample number from 0 to 1.0
func (s *SamplerLogger) SetSample(sample float32) {
	s.Sample = sample
}

// SetLevel set Logger level
// level can be log.AllLevel,log.debugLevel,...,log.OffLevel
func (s *SamplerLogger) SetLevel(level int) {
	s.LogLevel = level
}

// TraceField do a trace log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) TraceField(message field.Field, typ string, l field.InternalSpan) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", TraceLevelString)

	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))

	l.Record(r)
}

// DebugField do a debug log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) DebugField(message field.Field, typ string, l field.InternalSpan) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", DebugLevelString)
	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

// InfoField do a Info log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) InfoField(message field.Field, typ string, l field.InternalSpan) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", InfoLevelString)
	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

// WarnField do a Warn log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) WarnField(message field.Field, typ string, l field.InternalSpan) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", WarnLevelString)
	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

// ErrorField do a Error log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) ErrorField(message field.Field, typ string, l field.InternalSpan) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", ErrorLevelString)
	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

// FatalField do a Fatal log a object into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) FatalField(message field.Field, typ string, l field.InternalSpan) {
	if FatalLevel < s.LogLevel {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", FatalLevelString)
	r.Set(typ, message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

// Trace do a trace string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Trace(message string, l field.InternalSpan) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", TraceLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// Debug do a Debug string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Debug(message string, l field.InternalSpan) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", DebugLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// Info do a Info string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Info(message string, l field.InternalSpan) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", InfoLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// Warn do a Warn string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Warn(message string, l field.InternalSpan) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", WarnLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// Error do a Error string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Error(message string, l field.InternalSpan) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", ErrorLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// Fatal do a Fatal string log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) Fatal(message string, l field.InternalSpan) {
	if FatalLevel < s.LogLevel {
		return
	}
	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("SeverityText", FatalLevelString)
	r.Set("message", field.StringField(message))
	r.Set("timestamp", field.TimeField(time.Now()))
	l.Record(r)
}

// RecordMetrics do a metric log into InternalSpan,
// if InternalSpan is not nil, this interface will log the info,
// but not signal the InternalSpan
// if InternalSpan is nil, this interface will create a InternalSpan
// to log the info and signal the InternalSpan.
func (s *SamplerLogger) RecordMetrics(m field.Mmetric, l field.InternalSpan) {
	if l == nil {
		l = s.getInternalSpan()
		if l == nil {
			return
		}
		defer l.Signal()
	}

	l.Metric(m)
}

// NewInternalSpan return a root internal span
func (s *SamplerLogger) NewInternalSpan() field.InternalSpan {
	if s.runtime == nil {
		return nil
	}
	res := s.runtime.Children()
	// s.SetTraceID(field.GenTraceID(), res)
	return res
}

// ChildrenInternalSpan return a child InternalSpan for given InternalSpan
// If the InternalSpan is nil, will return nil
func (s *SamplerLogger) ChildrenInternalSpan(span field.InternalSpan) field.InternalSpan {
	if span == nil {
		return nil
	}
	res := span.Children()
	return res
}

// NewExternalSpan return a ExternalSpan for record exteranl call info.
// The ExternalSpan will created from InternalSpan.
// if InternalSpan is nil will return error
func (s *SamplerLogger) NewExternalSpan(span field.InternalSpan) (*field.ExternalSpanField, error) {
	if span == nil {
		return nil, field.NilPointerError
	}
	return span.NewExternalSpan(), nil
}

// SetParentID Set ParentId for the root InternalSpan
// If the InternalSpan is nil, will do nothing
func (s *SamplerLogger) SetParentID(ID string, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetParentID(ID)
}

// SetTraceID Set TraceId for the root InternalSpan
// If the InternalSpan is nil, will do nothing
func (s *SamplerLogger) SetTraceID(ID string, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetTraceID(ID)
}

// SetAttributes Set attributes for a root internalSpan
func (s *SamplerLogger) SetAttributes(t string, attrs field.Field, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetAttributes(t, attrs)
}

// func (s *SamplerLogger) Signal(span field.InternalSpan) {
// 	span.Signal()
// }
