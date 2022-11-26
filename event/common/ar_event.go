package common

import (
	"go.opentelemetry.io/otel/trace"
	"time"
)

// AREvent 对外暴露的 event 接口。
type AREvent interface {
	// SetEventType 设置非空 EventType 。
	SetEventType(eventType string)
	// SetTime 设置 Time 。
	SetTime(time time.Time)
	// SetLevel 设置 level 。
	SetLevel(level ARLevel)
	// SetAttributes 设置 Attributes 。
	SetAttributes(kvs ...ARAttribute)
	// SetSubject 设置 Subject 。
	SetSubject(subject string)
	// SetLink 设置 link 。
	SetLink(link trace.SpanContext)
	// SetData 设置 Data 。
	SetData(data interface{})

	// GetEventID 返回 EventID 。
	GetEventID() string
	// GetEventType 返回 EventType 。
	GetEventType() string
	// GetTime 返回 Time 。
	GetTime() time.Time
	// GetLevel 返回 level 。
	GetLevel() ARLevel
	// GetResource 返回 resource 。
	GetResource() ARResource
	// GetSubject 返回 Subject 。
	GetSubject() string
	// GetLink 返回 link 。
	GetLink() ARLink
	// GetData 返回 Data 。
	GetData() interface{}

	// GetEventMap 返回 map[string]interface{} 形式的 event 。
	GetEventMap() map[string]interface{}

	// private 禁止自己实现接口
	private()
	// setEventID 当前不允许修改 EventID 。
	// setEventID(eventID string)
}
