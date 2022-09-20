package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Link 自定义 Link 统一Attribute。
type Link struct {
	SpanContext           trace.SpanContext
	Attributes            []*Attribute
	DroppedAttributeCount int
}

// AnyRobotLinkFromLink 单条sdktrace.Link转换为*Link。
func AnyRobotLinkFromLink(link sdktrace.Link) *Link {
	return &Link{
		SpanContext:           link.SpanContext,
		Attributes:            AnyRobotAttributesFromKeyValues(link.Attributes),
		DroppedAttributeCount: link.DroppedAttributeCount,
	}
}

// AnyRobotLinksFromLinks 批量sdktrace.Link转换为[]*Link。
func AnyRobotLinksFromLinks(sdkLinks []sdktrace.Link) []*Link {
	if sdkLinks == nil {
		return make([]*Link, 0, 0)
	}
	links := make([]*Link, 0, len(sdkLinks))
	for i := 0; i < len(sdkLinks); i++ {
		links = append(links, AnyRobotLinkFromLink(sdkLinks[i]))
	}
	return links
}
