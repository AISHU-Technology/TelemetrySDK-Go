package eventsdk

import (
	"go.opentelemetry.io/otel/trace"
	"time"
)

// eventStartOptionFunc 执行 EventStartOption 的方法。
type eventStartOptionFunc func(*eventStartConfig) *eventStartConfig

func (o eventStartOptionFunc) apply(cfg *eventStartConfig) *eventStartConfig {
	return o(cfg)
}

func WithEventID(eventID string) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		if len(eventID) == 26 {
			cfg.EventID = eventID
		}
		return cfg
	})
}

func WithEventType(eventType string) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		if eventType != "" {
			cfg.EventType = eventType
		}
		return cfg
	})
}

func WithTime(t time.Time) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		if t.After(time.Time{}) {
			cfg.Time = t
		}
		return cfg
	})
}

func WithLevel(level Level) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Level = level
		return cfg
	})
}

func WithAttributes(attrs ...Attribute) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Attributes = attrs
		return cfg
	})
}

func WithSubject(subject string) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Subject = subject
		return cfg
	})
}

func WithSpanContext(spanContext trace.SpanContext) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.SpanContext = spanContext
		return cfg
	})
}

func WithData(data interface{}) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Data = data
		return cfg
	})
}
