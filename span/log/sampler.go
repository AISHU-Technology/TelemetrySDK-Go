// Provide some interfaces of logger
//
// looger can use to record string log of any object Log in the form of field
//
// logger can use to log internal thread/task with InteranlSpan.
//
// A LogSpan use to describe a thread job info or a task info. You can get a thread context info ever recorded.
// And then aggregation Span by TraceID or describe a Trace tree by SpanID and TraceID
package log

import (
	"github.com/hashicorp/go.net/context"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
	"math/rand"
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
	ctx      context.Context
}

func NewDefaultSamplerLogger() *SamplerLogger {
	return &SamplerLogger{
		Sample:   1.0,
		LogLevel: InfoLevel,
	}
}

func newRecord(typ string, message field.Field) field.Field {
	record := field.MallocStructField(2)
	record.Set(typ, message)
	record.Set("Type", field.StringField(typ))
	return record
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

func (s *SamplerLogger) getLogSpan() field.LogSpan {
	if s.runtime != nil {
		return s.runtime.Children(s.ctx)
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

// set context,if use trace,logSpan will Inheritance trace context,else ctx is nil or background
func (s *SamplerLogger) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *SamplerLogger) writeLogField(typ string, message, level field.Field, attr *field.Attribute) {
	l := s.getLogSpan()
	if l == nil {
		return
	}
	defer l.Signal()

	l.SetLogLevel(level)
	record := newRecord(typ, message)
	l.SetRecord(record)
	if attr != nil {
		l.SetAttributes(attr)
	}
}

func (s *SamplerLogger) writeLog(message string, level field.Field, attr *field.Attribute) {
	l := s.getLogSpan()
	if l == nil {
		return
	}
	defer l.Signal()

	l.SetLogLevel(level)
	record := field.MallocStructField(1)
	record.Set("Message", field.StringField(message))
	l.SetRecord(record)
	if attr != nil {
		l.SetAttributes(attr)
	}
}

// TraceField do a trace log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) TraceField(message field.Field, typ string, attr *field.Attribute) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLogField(typ, message, TraceLevelString, attr)

}

// DebugField do a debug log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) DebugField(message field.Field, typ string, attr *field.Attribute) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, DebugLevelString, attr)
}

// InfoField do a Info log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) InfoField(message field.Field, typ string, attr *field.Attribute) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLogField(typ, message, TraceLevelString, attr)
}

// WarnField do a Warn log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) WarnField(message field.Field, typ string, attr *field.Attribute) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, TraceLevelString, attr)
}

// ErrorField do a Error log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) ErrorField(message field.Field, typ string, attr *field.Attribute) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, TraceLevelString, attr)
}

// FatalField do a Fatal log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) FatalField(message field.Field, typ string, attr *field.Attribute) {
	if FatalLevel < s.LogLevel {
		return
	}

	s.writeLogField(typ, message, TraceLevelString, attr)
}

// Trace do a trace string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Trace(message string, attr *field.Attribute) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLog(message, TraceLevelString, attr)
}

// Debug do a Debug string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Debug(message string, attr *field.Attribute) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, DebugLevelString, attr)
}

// Info do a Info string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Info(message string, attr *field.Attribute) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, InfoLevelString, attr)
}

// Warn do a Warn string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Warn(message string, attr *field.Attribute) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, WarnLevelString, attr)
}

// Error do a Error string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Error(message string, attr *field.Attribute) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, ErrorLevelString, attr)
}

// Fatal do a Fatal string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Fatal(message string, attr *field.Attribute) {
	if FatalLevel < s.LogLevel {
		return
	}
	s.writeLog(message, FatalLevelString, attr)
}

//// NewLogSpan return a root internal span
//func (s *SamplerLogger) NewLogSpan(ctx context.Context) field.LogSpan {
//	if s.runtime == nil {
//		return nil
//	}
//	res := s.runtime.Children(ctx)
//
//	return res
//}

// SetAttributes Set attributes for a root LogSpan
//func (s *SamplerLogger) SetAttributes(t string, attrs field.Field, span field.LogSpan) {
//	if span == nil {
//		return
//	}
//
//	span.SetAttributes(t, attrs)
//}
