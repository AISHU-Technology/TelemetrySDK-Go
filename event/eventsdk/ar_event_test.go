package eventsdk

import (
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			"",
			args{""},
			NewEvent(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvent(tt.args.eventType); !reflect.DeepEqual(got.GetEventType(), tt.want.GetEventType()) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalEvents(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalEvents(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetData(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventID(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventID(); got != tt.want {
				t.Errorf("GetEventID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventMap(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventType(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventType(); got != tt.want {
				t.Errorf("GetEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetLevel(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Level
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetLink(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Link
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetLink(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetResource(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetSubject(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetSubject(); got != tt.want {
				t.Errorf("GetSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetTime(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventSend(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.Send()
		})
	}
}

func TestEventSetAttributes(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		kvs []Attribute
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetAttributes(tt.args.kvs...)
		})
	}
}

func TestEventSetData(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetData(tt.args.data)
		})
	}
}

func TestEventSetEventType(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		eventType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetEventType(tt.args.eventType)
		})
	}
}

func TestEventSetLevel(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		level Level
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetLevel(tt.args.level)
		})
	}
}

func TestEventSetLink(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		link trace.SpanContext
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetLink(tt.args.link)
		})
	}
}

func TestEventSetSubject(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetSubject(tt.args.subject)
		})
	}
}

func TestEventSetTime(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	type args struct {
		time time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetTime(tt.args.time)
		})
	}
}

func TestEventPrivate(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.private()
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateID(); got != tt.want {
				t.Errorf("generateID() = %v, want %v", got, tt.want)
			}
		})
	}
}
