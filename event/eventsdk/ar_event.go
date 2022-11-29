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
	Level     Level       `json:"Level"`
	Resource  *resource   `json:"Resource"`
	Subject   string      `json:"Subject"`
	Link      link        `json:"Link"`
	Data      interface{} `json:"Data"`
}

func Info(data interface{}, opts ...EventStartOption) Event {
	opts = append(opts, WithData(data))
	opts = append(opts, WithLevel(INFO))
	return NewEvent(opts...)
}

func Warn(data interface{}, opts ...EventStartOption) Event {
	opts = append(opts, WithData(data))
	opts = append(opts, WithLevel(WARN))
	return NewEvent(opts...)
}

func Error(data interface{}, opts ...EventStartOption) Event {
	opts = append(opts, WithData(data))
	opts = append(opts, WithLevel(ERROR))
	return NewEvent(opts...)
}

// NewEvent 创建新的 event ，默认填充ID、时间、事件级别、资源信息、事件类型。
func NewEvent(opts ...EventStartOption) Event {
	cfg := newEventStartConfig(opts...)
	e := &event{
		EventID:   cfg.EventID,
		EventType: cfg.EventType,
		Time:      cfg.Time,
		Level:     cfg.Level,
		Resource:  newResource(),
		Subject:   cfg.Subject,
		Link:      newLink(cfg.SpanContext),
		Data:      cfg.Data,
	}
	e.SetAttributes(cfg.Attributes...)
	return e
}

// generateID 生成全球唯一ULID。
func generateID() string {
	return ulid.Make().String()
}

func (e *event) SetEventID(eventID string) {
	if len(eventID) != 26 {
		log.Println(errors.New(custom_errors.ModuleName))
		return
	}
	e.EventID = eventID
}

func (e *event) SetEventType(eventType string) {
	if eventType == "" {
		log.Println(errors.New(custom_errors.ModuleName))
		return
	}
	e.EventType = eventType
}

func (e *event) SetTime(t time.Time) {
	if t.Equal(time.Time{}) {
		log.Println(errors.New(custom_errors.ModuleName))
		return
	}
	e.Time = t
}

func (e *event) SetLevel(level Level) {
	e.Level = level
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

func (e *event) SetLink(spanContext trace.SpanContext) {
	if !spanContext.IsValid() {
		log.Println(errors.New(custom_errors.ModuleName))
		return
	}
	e.Link.TraceID = spanContext.TraceID().String()
	e.Link.SpanID = spanContext.SpanID().String()
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
