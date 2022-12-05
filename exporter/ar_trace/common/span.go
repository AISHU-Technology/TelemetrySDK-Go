package common

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
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
	Resource             *Resource             `json:"Resource"`
	DroppedAttributes    int                   `json:"DroppedAttributes"`
	DroppedEvents        int                   `json:"DroppedEvents"`
	DroppedLinks         int                   `json:"DroppedLinks"`
	ChildSpanCount       int                   `json:"ChildSpanCount"`
}

// AnyRobotSpanFromReadOnlySpan 单条span转换为*AnyRobotSpan。
func AnyRobotSpanFromReadOnlySpan(span sdktrace.ReadOnlySpan) *AnyRobotSpan {
	if span == nil {
		return &AnyRobotSpan{}
	}
	return &AnyRobotSpan{Name: span.Name(), SpanContext: span.SpanContext(), Parent: span.Parent(), SpanKind: span.SpanKind(), StartTime: span.StartTime(), EndTime: span.EndTime(), Attributes: AnyRobotAttributesFromKeyValues(span.Attributes()), Links: AnyRobotLinksFromLinks(span.Links()), Events: AnyRobotEventsFromEvents(span.Events()), Status: span.Status(), InstrumentationScope: span.InstrumentationScope(), Resource: AnyRobotResourceFromResource(span.Resource()), DroppedAttributes: span.DroppedAttributes(), DroppedEvents: span.DroppedEvents(), DroppedLinks: span.DroppedLinks(), ChildSpanCount: span.ChildSpanCount()}
}

// AnyRobotSpansFromReadOnlySpans 批量span转换为[]*AnyRobotSpan。
func AnyRobotSpansFromReadOnlySpans(spans []sdktrace.ReadOnlySpan) []*AnyRobotSpan {
	if spans == nil {
		return make([]*AnyRobotSpan, 0)
	}
	arSpans := make([]*AnyRobotSpan, 0, len(spans))
	for _, span := range spans {
		arSpans = append(arSpans, AnyRobotSpanFromReadOnlySpan(span))
	}
	return arSpans
}
