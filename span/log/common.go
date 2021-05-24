package log

import (
	"span/field"
	"span/runtime"
)

func simpleDoc() {}

func allDoc() {}

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

	// NewInternalSpan return a root internal span
	NewInternalSpan() field.InternalSpan

	// SetParentID Set ParentId for the root InternalSpan
	// If the InternalSpan is nil, will do nothing
	SetParentID(ID string, span field.InternalSpan)

	// SetTraceID Set TraceId for the root InternalSpan
	// If the InternalSpan is nil, will do nothing
	SetTraceID(ID string, span field.InternalSpan)

	// NewExternalSpan return a ExternalSpan for record exteranl call info.
	// The ExternalSpan will created from InternalSpan.
	// if InternalSpan is nil will return error
	NewExternalSpan(span field.InternalSpan) (*field.ExternalSpanField, error)

	// ChildrenInternalSpan return a child InternalSpan for given InternalSpan
	// If the InternalSpan is nil, will return nil
	ChildrenInternalSpan(span field.InternalSpan) field.InternalSpan

	// SetAttributes Set attributes for a root internalSpan
	SetAttributes(t string, attrs field.Field, span field.InternalSpan)

	// TraceField do a trace log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	TraceField(message field.Field, typ string, l field.InternalSpan)

	// DebugField do a debug log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	DebugField(message field.Field, typ string, l field.InternalSpan)

	// InfoField do a Info log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	InfoField(message field.Field, typ string, l field.InternalSpan)

	// WarnField do a Warn log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	WarnField(message field.Field, typ string, l field.InternalSpan)

	// ErrorField do a Error log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	ErrorField(message field.Field, typ string, l field.InternalSpan)

	// FatalField do a Fatal log a object into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	FatalField(message field.Field, typ string, l field.InternalSpan)

	// Trace do a trace string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Trace(message string, l field.InternalSpan)

	// Debug do a Debug string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Debug(message string, l field.InternalSpan)

	// Info do a Info string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Info(message string, l field.InternalSpan)

	// Warn do a Warn string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Warn(message string, l field.InternalSpan)

	// Error do a Error string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Error(message string, l field.InternalSpan)

	// Fatal do a Fatal string log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	Fatal(message string, l field.InternalSpan)

	// RecordMetrics do a metric log into InternalSpan,
	// if InternalSpan is not nil, this interface will log the info,
	// but not signal the InternalSpan
	// if InternalSpan is nil, this interface will create a InternalSpan
	// to log the info and signal the InternalSpan.
	RecordMetrics(m field.Mmetric, l field.InternalSpan)
}
