package log

import (
	"os"
	"span/encoder"
	"span/field"
	"span/log"
	"span/open_standard"
	"span/runtime"
)

var _defaultSamplelogger *log.SamplerLogger

func init() {
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(os.Stdout),
	}
	deafaultRuntime := runtime.NewRuntime(writer, field.NewSpanFromPool)
	_defaultSamplelogger = log.NewdefaultSamplerLogger()
	_defaultSamplelogger.SetRuntime(deafaultRuntime)
	// _defaultSamplelogger.Runtime = deafaultRuntime
	go deafaultRuntime.Run()
}

// SetLogger() will set default backend Samplerlogger as given after the old Samplerlogger is closed
// but SetLogger() will not start the Runtime, caller should start by themself
func SetLogger(l *log.SamplerLogger) {
	_defaultSamplelogger.Close()
	_defaultSamplelogger = l
}

func Close() {
	_defaultSamplelogger.Close()
}

func TraceField(message field.Field, l field.InternalSpan) {
	_defaultSamplelogger.TraceField(message, l)
}

func DebugField(message field.Field, l field.InternalSpan) {
	_defaultSamplelogger.TraceField(message, l)
}

func WarnField(message field.Field, l field.InternalSpan) {
	_defaultSamplelogger.TraceField(message, l)
}

func ErrorField(message field.Field, l field.InternalSpan) {
	_defaultSamplelogger.TraceField(message, l)
}

func FatalField(message field.Field, l field.InternalSpan) {
	_defaultSamplelogger.TraceField(message, l)
}

func Trace(message string, l field.InternalSpan) {
	_defaultSamplelogger.Trace(message, l)
}

func Debug(message string, l field.InternalSpan) {
	_defaultSamplelogger.Debug(message, l)
}

func Warn(message string, l field.InternalSpan) {
	_defaultSamplelogger.Warn(message, l)
}

func Error(message string, l field.InternalSpan) {
	_defaultSamplelogger.Error(message, l)
}

func Fatal(message string, l field.InternalSpan) {
	_defaultSamplelogger.Fatal(message, l)
}

func RecordMetrics(m field.Mmetric, s field.InternalSpan) {
	_defaultSamplelogger.RecordMetrics(m, s)
}
