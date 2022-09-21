package common

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"time"
)

// AnyRobotSpan 实现ReadOnlySpan接口的结构体，属性无增删。
type AnyRobotSpan struct {
	Name                 string                `json:"Name"`
	SpanContext          trace.SpanContext     `json:"SpanContext"`
	Parent               trace.SpanContext     `json:"Parent"`
	SpanKind             trace.SpanKind        `json:"SpanKind"`
	StartTime            time.Time             `json:"StartTime"`
	EndTime              time.Time             `json:"EndTime"`
	Attributes           []*Attribute          `json:"Attributes"`
	Links                []*Link               `json:"Links"`
	Events               []*Event              `json:"Events"`
	Status               sdktrace.Status       `json:"Status"`
	InstrumentationScope instrumentation.Scope `json:"InstrumentationScope"`
	Resource             *resource.Resource    `json:"Resource"`
	DroppedAttributes    int                   `json:"DroppedAttributes"`
	DroppedEvents        int                   `json:"DroppedEvents"`
	DroppedLinks         int                   `json:"DroppedLinks"`
	ChildSpanCount       int                   `json:"ChildSpanCount"`
}

// AnyRobotSpanFromReadOnlySpan 单条span转换为*AnyRobotSpan。
func AnyRobotSpanFromReadOnlySpan(ros sdktrace.ReadOnlySpan) *AnyRobotSpan {
	if ros == nil {
		return &AnyRobotSpan{}
	}
	return &AnyRobotSpan{
		Name:                 ros.Name(),
		SpanContext:          ros.SpanContext(),
		Parent:               ros.Parent(),
		SpanKind:             ros.SpanKind(),
		StartTime:            ros.StartTime(),
		EndTime:              ros.EndTime(),
		Attributes:           AnyRobotAttributesFromKeyValues(ros.Attributes()),
		Links:                AnyRobotLinksFromLinks(ros.Links()),
		Events:               AnyRobotEventsFromEvents(ros.Events()),
		Status:               ros.Status(),
		InstrumentationScope: ros.InstrumentationScope(),
		Resource:             ros.Resource(),
		DroppedAttributes:    ros.DroppedAttributes(),
		DroppedEvents:        ros.DroppedEvents(),
		DroppedLinks:         ros.DroppedLinks(),
		ChildSpanCount:       ros.ChildSpanCount(),
	}
}

// AnyRobotSpansFromReadOnlySpans 批量span转换为[]*AnyRobotSpan。
func AnyRobotSpansFromReadOnlySpans(ross []sdktrace.ReadOnlySpan) []*AnyRobotSpan {
	if ross == nil {
		return make([]*AnyRobotSpan, 0)
	}
	spans := make([]*AnyRobotSpan, 0, len(ross))
	for i := 0; i < len(ross); i++ {
		spans = append(spans, AnyRobotSpanFromReadOnlySpan(ross[i]))
	}
	return spans
}
