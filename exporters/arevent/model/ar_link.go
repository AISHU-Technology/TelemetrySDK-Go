package model

// ARLink 和 Trace 关联，记录 TraceID 和 SpanID 。
type ARLink interface {
	// GetTraceID 返回 TraceID 。
	GetTraceID() string
	// GetSpanID 返回 SpanID 。
	GetSpanID() string
}
