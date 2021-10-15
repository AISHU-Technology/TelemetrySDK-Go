package log

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
)

// log level text format
const (
	TraceLevelString = field.StringField("Trace")
	DebugLevelString = field.StringField("Debug")
	InfoLevelString  = field.StringField("Info")
	WarnLevelString  = field.StringField("Warn")
	ErrorLevelString = field.StringField("Error")
	FatalLevelString = field.StringField("Fatal")
)

// log level int format
const (
	AllLevel = iota
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	OffLevel
)

type Logger interface {
	// SetSample set the log sample
	// sample number from 0 to 1.0
	SetSample(sample float32)

	// SetLevel set Logger level
	// level can be log.AllLevel,log.debugLevel,...,log.OffLevel
	SetLevel(level int)

	// close logger
	Close()

	// SetRuntime for logger
	SetRuntime(*runtime.Runtime)
	// TraceField do a trace log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	TraceField(message field.Field, typ string, options ...field.LogOptionFunc)

	// DebugField do a debug log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	DebugField(message field.Field, typ string, options ...field.LogOptionFunc)

	// InfoField do a Info log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	InfoField(message field.Field, typ string, options ...field.LogOptionFunc)

	// WarnField do a Warn log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	WarnField(message field.Field, typ string, options ...field.LogOptionFunc)

	// ErrorField do a Error log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	ErrorField(message field.Field, typ string, options ...field.LogOptionFunc)

	// FatalField do a Fatal log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	FatalField(message field.Field, typ string, options ...field.LogOptionFunc)

	// Trace do a trace string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Trace(message string, options ...field.LogOptionFunc)

	// Debug do a Debug string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// but not signal the LogSpan
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Debug(message string, options ...field.LogOptionFunc)

	// Info do a Info string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Info(message string, options ...field.LogOptionFunc)

	// Warn do a Warn string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Warn(message string, options ...field.LogOptionFunc)

	// Error do a Error string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Error(message string, options ...field.LogOptionFunc)

	// Fatal do a Fatal string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Fatal(message string, options ...field.LogOptionFunc)
}
