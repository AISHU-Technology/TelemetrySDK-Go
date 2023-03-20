package log

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
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
