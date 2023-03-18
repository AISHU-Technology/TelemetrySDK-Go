package log

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/runtime"
	"math/rand"
	"time"
)

// SyncLogger 同步发送模式的日志器。
type SyncLogger interface {
	SetSample(sample float32)
	SetLogLevel(logLevel int)
	TraceField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	DebugField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	InfoField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	WarnField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	ErrorField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	FatalField(message field.Field, type_ string, opts ...field.LogOptionFunc) error
	Trace(message string, opts ...field.LogOptionFunc) error
	Debug(message string, opts ...field.LogOptionFunc) error
	Info(message string, opts ...field.LogOptionFunc) error
	Warn(message string, opts ...field.LogOptionFunc) error
	Error(message string, opts ...field.LogOptionFunc) error
	Fatal(message string, opts ...field.LogOptionFunc) error
}

// syncLogger 同步发送模式的日志器。
type syncLogger struct {
	Sample   float32
	LogLevel int
	runtime  *runtime.Runtime
	ctx      context.Context
}

// NewSyncLogger 创建同步发送模式的日志器，创建时可传入参数。记日志的方法有返回值，返回error=nil代表发送成功，返回error!=nil代表发送失败。
func NewSyncLogger(opts ...LoggerStartOption) SyncLogger {
	cfg := newLoggerStartConfig(opts...)
	return &syncLogger{
		Sample:   cfg.Sample,
		LogLevel: cfg.LogLevel,
	}
}

// SetLogLevel 设置日志级别，从0~7，0代表全部输出，7代表关闭输出。
func (logger *syncLogger) SetLogLevel(logLevel int) {
	if logLevel < AllLevel || logLevel > OffLevel {
		return
	}
	logger.LogLevel = logLevel
}

// SetSample 设置采样等级，从0.0~1.0，0.0代表不采样，1.0代表全采样。
func (logger *syncLogger) SetSample(sample float32) {
	if sample < 0.0 || sample > 1.0 {
		return
	}
}

// TraceField Trace 级别的日志，记录结构体。
func (logger *syncLogger) TraceField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if TraceLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(type_, message, TraceLevelString, opts...)
}

// DebugField Debug 级别的日志，记录结构体。
func (logger *syncLogger) DebugField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if DebugLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(type_, message, DebugLevelString, opts...)
}

// InfoField Info 级别的日志，记录结构体。
func (logger *syncLogger) InfoField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if InfoLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(type_, message, InfoLevelString, opts...)
}

// WarnField Warn 级别的日志，记录结构体。
func (logger *syncLogger) WarnField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if WarnLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(type_, message, WarnLevelString, opts...)
}

// ErrorField Error 级别的日志，记录结构体。
func (logger *syncLogger) ErrorField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if ErrorLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(type_, message, ErrorLevelString, opts...)
}

// FatalField Fatal 级别的日志，记录结构体。
func (logger *syncLogger) FatalField(message field.Field, type_ string, opts ...field.LogOptionFunc) error {
	if FatalLevel < logger.LogLevel {
		return nil
	}
	return logger.writeLogField(type_, message, FatalLevelString, opts...)
}

// Trace Trace 级别的日志，记录字符串。
func (logger *syncLogger) Trace(message string, opts ...field.LogOptionFunc) error {
	if TraceLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, TraceLevelString, opts...)
}

// Debug Debug 级别的日志，记录字符串。
func (logger *syncLogger) Debug(message string, opts ...field.LogOptionFunc) error {
	if DebugLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, DebugLevelString, opts...)
}

// Info Info 级别的日志，记录字符串。
func (logger *syncLogger) Info(message string, opts ...field.LogOptionFunc) error {
	if InfoLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, InfoLevelString, opts...)
}

// Warn Warn 级别的日志，记录字符串。
func (logger *syncLogger) Warn(message string, opts ...field.LogOptionFunc) error {
	if WarnLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, WarnLevelString, opts...)
}

// Error Error 级别的日志，记录字符串。
func (logger *syncLogger) Error(message string, opts ...field.LogOptionFunc) error {
	if ErrorLevel < logger.LogLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, ErrorLevelString, opts...)
}

// Fatal Fatal 级别的日志，记录字符串。
func (logger *syncLogger) Fatal(message string, opts ...field.LogOptionFunc) error {
	if FatalLevel < logger.LogLevel {
		return nil
	}
	return logger.writeLog(message, FatalLevelString, opts...)
}

// sampleCheck 检查采样率决定是否记录当前日志。
func (logger *syncLogger) sampleCheck() bool {
	// 全采样
	if logger.Sample >= 1.0 {
		return true
	}
	// 全丢弃
	if logger.Sample <= 0 {
		return false
	}
	// 生成0.0~1.0之间的随机数
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Float32() <= logger.Sample
}

// getLogSpan 获取Log上下文。
func (logger *syncLogger) getLogSpan() field.LogSpan {
	if logger.runtime != nil {
		return logger.runtime.Children(logger.ctx)
	}
	return nil
}

// writeLogField 写非结构化日志。！！！！！！！！！！什么时候发送日志，并且返回错误？
func (logger *syncLogger) writeLogField(typ string, message, level field.Field, options ...field.LogOptionFunc) error {
	logSpan := logger.getLogSpan()
	if logSpan == nil {
		return nil
	}
	defer logSpan.Signal()
	logSpan.SetLogLevel(level)
	record := newRecord(typ, message)
	logSpan.SetRecord(record)
	logSpan.SetOption(options...)
	return nil
}

// writeLog 写结构化日志。！！！！！！！！！！什么时候发送日志，并且返回错误？
func (logger *syncLogger) writeLog(message string, level field.Field, options ...field.LogOptionFunc) error {
	logSpan := logger.getLogSpan()
	if logSpan == nil {
		return nil
	}
	defer logSpan.Signal()
	logSpan.SetLogLevel(level)
	record := field.MallocStructField(1)
	record.Set("Message", field.StringField(message))
	logSpan.SetRecord(record)
	logSpan.SetOption(options...)
	return nil
}
