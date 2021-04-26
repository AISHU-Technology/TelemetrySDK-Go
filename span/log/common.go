package log

import (
	"span/field"
)

var (
	TraceLevelString = field.StringField("Trace")
	DebugLevelString = field.StringField("Debug")
	InfoLevelString  = field.StringField("Info")
	WarnLevelString  = field.StringField("Warn")
	ErrorLevelString = field.StringField("Error")
	FatalLevelString = field.StringField("Fatal")
)

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
