package common

// Link 用于关联 Trace 信息， Event 和 Trace 一对一。
type Link struct {
	TraceID string `json:"TraceID"`
	SpanID  string `json:"SpanID"`
}

// NewLink 创建新的 Link 。
func NewLink() Link {
	return Link{}
}

func (l Link) GetTraceID() string {
	return l.TraceID
}

func (l Link) GetSpanID() string {
	return l.SpanID
}
