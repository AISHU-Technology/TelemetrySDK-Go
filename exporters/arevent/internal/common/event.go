package common

import (
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/errors"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type Event struct {
	EventID   string        `json:"EventID"`
	EventType string        `json:"EventType"`
	Time      time.Time     `json:"Time"`
	Level     model.ARLevel `json:"Level"`
	Resource  *Resource     `json:"Resource"`
	Subject   string        `json:"Subject"`
	Link      Link          `json:"Link"`
	Data      interface{}   `json:"Data"`
}

// DefaultEventType 默认的非空事件类型
const DefaultEventType = "Telemetry.Default.Event"

// NewEvent 创建新的 Event ，默认填充ID、时间、事件级别、资源信息，需要传入事件类型，默认为"Telemetry.Default.Event"。
func NewEvent(eventType string) model.AREvent {
	if eventType == "" {
		eventType = DefaultEventType
	}
	return &Event{
		EventID:   generateID(),
		EventType: eventType,
		Time:      time.Now(),
		Level:     INFO,
		Resource:  NewResource(),
		Subject:   "",
		Link:      NewLink(),
		Data:      nil,
	}
}

// generateID 生成全球唯一ULID。
func generateID() string {
	return ulid.Make().String()
}

// setEventID 当前版本不允许设置 EventID 。
//func (e *Event) setEventID(eventID string) {
//	e.EventID = eventID
//}

// SetEventType 设置非空 EventType 。
func (e *Event) SetEventType(eventType string) {
	if eventType == "" {
		eventType = DefaultEventType
	}
	e.EventType = eventType
}

// SetTime 设置 time.Time 类型的 Time 。
func (e *Event) SetTime(time time.Time) {
	e.Time = time
}

// SetLevel 设置事件级别 ARLevel 。
func (e *Event) SetLevel(level model.ARLevel) {
	e.Level = level
}

// SetAttributes 设置资源 Resource 。
func (e *Event) SetAttributes(kvs ...model.ARAttribute) {
	for _, kv := range kvs {
		// 校验 Attribute 是否合法，合法的才放进map去重。
		if !kv.Valid() {
			log.Println(customErrors.AnyRobotEventExporter_EmptyKey)
			continue
		}
		e.Resource.AttributesMap[kv.GetKey()] = kv.GetValue()
	}
	// 去重map转数组。
	e.Resource.mapToSlice()
}

// SetSubject 设置操作对象 Subject 。
func (e *Event) SetSubject(subject string) {
	e.Subject = subject
}

// SetLink 设置关联链路 Link 。
func (e *Event) SetLink(link trace.SpanContext) {
	e.Link.TraceID = link.TraceID()
	e.Link.SpanID = link.SpanID()
}

// SetData 设置事件数据 Data 。
func (e *Event) SetData(data interface{}) {
	e.Data = data
}

// GetEventID 获取事件唯一标识符 EventID 。
func (e *Event) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型 EventType 。
func (e *Event) GetEventType() string {
	return e.EventType
}

// GetTime 获取事件时间 Time 。
func (e *Event) GetTime() time.Time {
	return e.Time
}

// GetLevel 获取事件级别 Level 。
func (e *Event) GetLevel() model.ARLevel {
	return e.Level
}

// GetResource 获取事件资源 *Resource 。
func (e *Event) GetResource() model.ARResource {
	return e.Resource
}

// GetSubject 获取事件操作对象 Subject 。
func (e *Event) GetSubject() string {
	return e.Subject
}

// GetLink 获取关联链路ID Link.TraceID Link.SpanID 。
func (e *Event) GetLink() model.ARLink {
	return e.Link
}

// GetData 获取事件数据 Data 。
func (e *Event) GetData() interface{} {
	return e.Data
}

// private 禁止实现 AREvent 接口。
//func (e *Event) private() {}
