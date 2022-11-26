package common

import (
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/errors"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

// Event 自定义 Event 统一数据模型。
type Event struct {
	EventID   string      `json:"EventID"`
	EventType string      `json:"EventType"`
	Time      time.Time   `json:"Time"`
	Level     Level       `json:"Level"`
	Resource  *Resource   `json:"Resource"`
	Subject   string      `json:"Subject"`
	Link      Link        `json:"Link"`
	Data      interface{} `json:"Data"`
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

func (e *Event) SetEventType(eventType string) {
	if eventType == "" {
		eventType = DefaultEventType
	}
	e.EventType = eventType
}

func (e *Event) SetTime(time time.Time) {
	e.Time = time
}

func (e *Event) SetLevel(level model.ARLevel) {
	e.Level = Level(level.Self())
}

func (e *Event) SetAttributes(kvs ...model.ARAttribute) {
	for _, kv := range kvs {
		// 校验 Attribute 是否合法，合法的才放进map去重。
		if !kv.Valid() {
			log.Println(customErrors.AnyRobotEventExporter_InvalidKey)
			continue
		}
		e.Resource.AttributesMap[kv.GetKey()] = kv.GetValue().GetData()
	}
}

func (e *Event) SetSubject(subject string) {
	e.Subject = subject
}

func (e *Event) SetLink(link trace.SpanContext) {
	e.Link.TraceID = link.TraceID().String()
	e.Link.SpanID = link.SpanID().String()
}

func (e *Event) SetData(data interface{}) {
	e.Data = data
}

func (e *Event) GetEventID() string {
	return e.EventID
}

func (e *Event) GetEventType() string {
	return e.EventType
}

func (e *Event) GetTime() time.Time {
	return e.Time
}

func (e *Event) GetLevel() model.ARLevel {
	return e.Level
}

func (e *Event) GetResource() model.ARResource {
	return e.Resource
}

func (e *Event) GetSubject() string {
	return e.Subject
}

func (e *Event) GetLink() model.ARLink {
	return e.Link
}

func (e *Event) GetData() interface{} {
	return e.Data
}

func (e *Event) GetEventMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["EventID"] = e.EventID
	result["EventType"] = e.EventType
	result["Time"] = e.Time
	result["Level"] = e.Level
	result["Resource"] = e.Resource
	result["Subject"] = e.Subject
	result["Link"] = e.Link
	result["Data"] = e.Data

	return result
}
