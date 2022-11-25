package model

import (
	"go.opentelemetry.io/otel/trace"
	"time"
)

type AREvent interface {
	// SetEventType 设置非空 EventType 。
	SetEventType(eventType string)
	// SetTime 设置 Time 。
	SetTime(time time.Time)
	// SetLevel 设置 Level 。
	SetLevel(level ARLevel)
	// SetAttributes 设置 Attributes 。
	SetAttributes(kvs ...ARAttribute)
	// SetSubject 设置 Subject 。
	SetSubject(subject string)
	// SetLink 设置 Link 。
	SetLink(link trace.SpanContext)
	// SetData 设置 Data 。
	SetData(data interface{})

	// GetEventID 返回 EventID 。
	GetEventID() string
	// GetEventType 返回 EventType 。
	GetEventType() string
	// GetTime 返回 Time 。
	GetTime() time.Time
	// GetLevel 返回 Level 。
	GetLevel() ARLevel
	// GetResource 返回 Resource 。
	GetResource() ARResource
	// GetSubject 返回 Subject 。
	GetSubject() string
	// GetLink 返回 Link 。
	GetLink() ARLink
	// GetData 返回 Data 。
	GetData() interface{}

	// GetEventMap 返回 map[string]interface{} 形式的 Event 。
	GetEventMap() map[string]interface{}

	// setEventID 当前不允许修改 EventID 。
	// setEventID(eventID string)
}
