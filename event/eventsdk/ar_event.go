package eventsdk

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/custom_errors"
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
	Level     level       `json:"Level"`
	Resource  *resource   `json:"Resource"`
	Subject   string      `json:"Subject"`
	Link      link        `json:"Link"`
	Data      interface{} `json:"Data"`
}

// DefaultEventType 默认的非空事件类型
const DefaultEventType = "Default.EventType"

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
			log.Println(custom_errors.Event_InvalidKey)
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
	result["Level"] = e.Level
	result["Resource"] = e.Resource
	result["Subject"] = e.Subject
	result["Link"] = e.Link
	result["Data"] = e.Data

	return result
}

func (e *event) Send() {
	GetEventProvider().(*eventProvider).loadEvent(e)
}

func (e *event) private() {}

// UnmarshalEvents 将JSON解析成[]Event。
func UnmarshalEvents(b []byte) ([]Event, error) {
	// 使用实际类型 []*event 来接收数据。
	events := make([]*event, 0)
	err := json.Unmarshal(b, &events)

	// 返回接口类型 Event 。
	result := make([]Event, 0, len(events))
	for _, e := range events {
		// 校验格式
		if e != nil && e.Valid() {
			result = append(result, e)
		}
	}
	// 如果返回空切片说明传入的JSON格式错误。
	if len(result) == 0 {
		err = errors.New(custom_errors.Event_InvalidJSON)
	}
	return result, err
}

func (e *event) Valid() bool {

	return len(e.GetEventID()) == 26 && e.GetEventType() != "" && e.GetTime().After(time.Time{}) && e.GetLevel().Valid() && e.GetResource().Valid() && e.GetLink().Valid()
}
