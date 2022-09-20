package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

// Event 自定义 Event 统一Attribute。
type Event struct {
	Name                  string
	Attributes            []*Attribute
	DroppedAttributeCount int
	Time                  time.Time
}

// AnyRobotEventFromEvent 单条sdktrace.Event转换为*Event。
func AnyRobotEventFromEvent(event sdktrace.Event) *Event {
	return &Event{
		Name:                  event.Name,
		Attributes:            AnyRobotAttributesFromKeyValues(event.Attributes),
		DroppedAttributeCount: event.DroppedAttributeCount,
		Time:                  event.Time,
	}
}

// AnyRobotEventsFromEvents 批量sdktrace.Event转换为[]*Event。
func AnyRobotEventsFromEvents(sdkEvent []sdktrace.Event) []*Event {
	if sdkEvent == nil {
		return make([]*Event, 0, 0)
	}
	events := make([]*Event, 0, len(sdkEvent))
	for i := 0; i < len(sdkEvent); i++ {
		events = append(events, AnyRobotEventFromEvent(sdkEvent[i]))
	}
	return events
}
