package eventsdk

// link 用于关联 Trace 信息， event 和 Trace 一对一。
type link struct {
	TraceID string `json:"TraceID"`
	SpanID  string `json:"SpanID"`
}

// newLink 创建新的 link 。
func newLink() link {
	return link{
		TraceID: "00000000000000000000000000000000",
		SpanID:  "0000000000000000",
	}
}

func (l link) GetTraceID() string {
	return l.TraceID
}

func (l link) GetSpanID() string {
	return l.SpanID
}

func (l link) Valid() bool {
	return len(l.GetTraceID()) == 32 && len(l.GetSpanID()) == 16
}

func (l link) private() {}
