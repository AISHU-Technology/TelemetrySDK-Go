package eventsdk

import (
	"go.opentelemetry.io/otel/trace"
	"time"
)

// Event 对外暴露的 event 接口。
type Event interface {
	// SetEventType 设置非空 EventType 。
	SetEventType(eventType string)
	// SetTime 设置 Time 。
	SetTime(time time.Time)
	//SetLevel 设置 level 。
	SetLevel(level Level)
	// SetAttributes 设置 Attributes 。
	SetAttributes(kvs ...Attribute)
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
	GetLevel() Level
	// GetResource 返回 resource 。
	GetResource() Resource
	// GetSubject 返回 Subject 。
	GetSubject() string
	// GetLink 返回 link 。
	GetLink() Link
	// GetData 返回 Data 。
	GetData() interface{}
	// GetEventMap 返回 map[string]interface{} 形式的 event 。
	GetEventMap() map[string]interface{}
	// Send 上报 Event 到 AnyRobot Event 数据接收器。
	Send()

	// private 禁止用户自己实现接口。
	private()
	//SetEventID 当前不允许修改 EventID 。
	//SetEventID(eventID string)

}
