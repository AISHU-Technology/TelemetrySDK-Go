package log

import (
	"math/rand"
	"span"
	"span/field"
	"time"
)

type SamplerLogger struct {
	// logger sample
	Sample float32
	// logger info
	LogLevel int
	Runtime  *span.Runtime
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
	s.Runtime.Signal()
}

func (s *SamplerLogger) sampleCheck() bool {
	if s.Sample >= 100 {
		return true
	}

	if s.Sample <= 0 {
		return false
	}

	rand.Seed(time.Now().UnixNano())

	return rand.Float32() <= s.Sample
}

func (s *SamplerLogger) RecordMetrics(m field.Mmetric, l field.InternalSpan) {
	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}

	l.Metric(m)
}

func (s *SamplerLogger) TraceField(message field.Field, l field.InternalSpan) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", TraceLevelString)

	r.Set("message", field.StringField("message...."))
}

func (s *SamplerLogger) DebugField(message field.Field, l field.InternalSpan) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", DebugLevelString)
	r.Set("message", message)
}

func (s *SamplerLogger) WarnField(message field.Field, l field.InternalSpan) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", WarnLevelString)
	r.Set("message", message)
}

func (s *SamplerLogger) ErrorField(message field.Field, l field.InternalSpan) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", ErrorLevelString)
	r.Set("message", message)
}

func (s *SamplerLogger) FatalField(message field.Field, l field.InternalSpan) {
	if FatalLevel < s.LogLevel {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", FatalLevelString)
	r.Set("message", message)
}

func (s *SamplerLogger) Trace(message string, l field.InternalSpan) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", TraceLevelString)
	r.Set("message", field.StringField("message..."))
}

func (s *SamplerLogger) Debug(message string, l field.InternalSpan) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", DebugLevelString)
	r.Set("message", field.StringField("message..."))
}

func (s *SamplerLogger) Warn(message string, l field.InternalSpan) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", WarnLevelString)
	r.Set("message", field.StringField("message..."))
}

func (s *SamplerLogger) Error(message string, l field.InternalSpan) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", ErrorLevelString)
	r.Set("message", field.StringField("message..."))
}

func (s *SamplerLogger) Fatal(message string, l field.InternalSpan) {
	if FatalLevel < s.LogLevel {
		return
	}
	if l == nil {
		l = s.Runtime.Children()
		if l == nil {
			return
		}
		defer l.Signal()
	}
	r := newStructRecord()
	r.Set("level", FatalLevelString)
	r.Set("message", field.StringField("message..."))
}
