package model

import (
	"go.opentelemetry.io/otel/trace"
	"time"
)

type AREvent interface {
	//setEventID(eventID string)
	SetEventType(eventType string)
	SetTime(time time.Time)
	SetLevel(level ARLevel)
	SetAttributes(kvs ...ARAttribute)
	SetSubject(subject string)
	SetLink(link trace.SpanContext)
	SetData(data interface{})

	GetEventID() string
	GetEventType() string
	GetTime() time.Time
	GetLevel() ARLevel
	GetResource() ARResource
	GetSubject() string
	GetLink() ARLink
	GetData() interface{}

	//private()
}
