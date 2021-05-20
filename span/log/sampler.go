package log

import (
	"math/rand"
	"span/field"
	"span/runtime"
	"time"
)

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

func (s *SamplerLogger) Close() {
	if s.runtime != nil {
		s.runtime.Signal()
	}
}

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

func (s *SamplerLogger) SetSample(sample float32) {
	s.Sample = sample
}

func (s *SamplerLogger) SetLevel(level int) {
	s.LogLevel = level
}

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

	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))

	l.Record(r)
}

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
	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

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
	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

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
	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

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
	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

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
	r.Set("message", message)
	r.Set("timestamp", field.TimeField(time.Now()))
	r.Set("type", field.StringField(typ))
	l.Record(r)
}

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

func (s *SamplerLogger) NewInternalSpan() field.InternalSpan {
	if s.runtime == nil {
		return nil
	}
	res := s.runtime.Children()
	// s.SetTraceID(field.GenTraceID(), res)
	return res
}

func (s *SamplerLogger) ChildrenInternalSpan(span field.InternalSpan) field.InternalSpan {
	if span == nil {
		return nil
	}
	res := span.Children()
	return res
}

func (s *SamplerLogger) NewExternalSpan(span field.InternalSpan) (*field.ExternalSpanField, error) {
	if span == nil {
		return nil, field.NilPointerError
	}
	return span.NewExternalSpan(), nil
}

func (s *SamplerLogger) SetParentID(ID string, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetParentID(ID)
}

func (s *SamplerLogger) SetTraceID(ID string, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetTraceID(ID)
}

func (s *SamplerLogger) SetAttributes(t string, attrs field.Field, span field.InternalSpan) {
	if span == nil {
		return
	}

	span.SetAttributes(t, attrs)
}

// func (s *SamplerLogger) Signal(span field.InternalSpan) {
// 	span.Signal()


