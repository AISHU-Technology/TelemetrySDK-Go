package eventsdk

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/customerrors"
	"encoding/json"
	"errors"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

// event 自定义 event 统一数据模型。
type event struct {
	EventID   string      `json:"EventID"`
	EventType string      `json:"EventType"`
	Time      time.Time   `json:"Time"`
	Level     level       `json:"level"`
	Resource  *resource   `json:"resource"`
	Subject   string      `json:"Subject"`
	Link      link        `json:"link"`
	Data      interface{} `json:"Data"`
}

// DefaultEventType 默认的非空事件类型
const DefaultEventType = "Telemetry.Default.Event"

// NewEvent 创建新的 event ，默认填充ID、时间、事件级别、资源信息，需要传入事件类型，默认为"Telemetry.Default.Event"。
func NewEvent(eventType string) Event {
	if eventType == "" {
		eventType = DefaultEventType
	}
	return &event{
		EventID:   generateID(),
		EventType: eventType,
		Time:      time.Now(),
		Level:     INFO,
		Resource:  newResource(),
		Subject:   "",
		Link:      newLink(),
		Data:      nil,
	}
}

// generateID 生成全球唯一ULID。
func generateID() string {
	return ulid.Make().String()
}

func (e *event) SetEventType(eventType string) {
	if eventType == "" {
		eventType = DefaultEventType
	}
	e.EventType = eventType
}

func (e *event) SetTime(time time.Time) {
	e.Time = time
}

func (e *event) SetLevel(level Level) {
	e.Level = newLevel(level.Self())
}

func (e *event) SetAttributes(kvs ...Attribute) {
	for _, kv := range kvs {
		// 校验 attribute 是否合法，合法的才放进map去重。
		if !kv.Valid() {
			log.Println(customerrors.Event_InvalidKey)
			continue
		}
		e.Resource.AttributesMap[kv.GetKey()] = kv.GetValue().GetData()
	}
}

func (e *event) SetSubject(subject string) {
	e.Subject = subject
}

func (e *event) SetLink(link trace.SpanContext) {
	e.Link.TraceID = link.TraceID().String()
	e.Link.SpanID = link.SpanID().String()
}

func (e *event) SetData(data interface{}) {
	e.Data = data
}

func (e *event) GetEventID() string {
	return e.EventID
}

func (e *event) GetEventType() string {
	return e.EventType
}

func (e *event) GetTime() time.Time {
	return e.Time
}

func (e *event) GetLevel() Level {
	return e.Level
}

func (e *event) GetResource() Resource {
	return e.Resource
}

func (e *event) GetSubject() string {
	return e.Subject
}

func (e *event) GetLink() Link {
	return e.Link
}

func (e *event) GetData() interface{} {
	return e.Data
}

func (e *event) GetEventMap() map[string]interface{} {
	result := make(map[string]interface{}, 8)
	result["EventID"] = e.EventID
	result["EventType"] = e.EventType
	result["Time"] = e.Time
	result["level"] = e.Level
	result["resource"] = e.Resource
	result["Subject"] = e.Subject
	result["link"] = e.Link
	result["Data"] = e.Data

	return result
}

func (e *event) Send() {
	GetEventProvider().(*eventProvider).loadEvent(e)
}

func (e *event) private() {}

// UnmarshalEvents 将JSON解析成[]Event。
func UnmarshalEvents(b []byte) ([]Event, error) {
	events := make([]*event, 0)
	err := json.Unmarshal(b, &events)

	result := make([]Event, 0, len(events))
	for _, e := range events {
		result = append(result, e)
	}
	if len(result) == 0 {
		err = errors.New(customerrors.Event_InvalidJSON)
	}
	return result, err
}
