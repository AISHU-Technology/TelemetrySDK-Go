package eventsdk

// Link 和 Trace 关联，记录 TraceID 和 SpanID 。
type Link interface {
	// GetTraceID 返回 TraceID 。
	GetTraceID() string
	// GetSpanID 返回 SpanID 。
	GetSpanID() string
	// Valid 校验是否合法。
	Valid() bool
	// private 禁止用户自己实现接口。
	private()
}
