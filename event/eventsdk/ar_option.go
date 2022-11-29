package eventsdk

import (
	"time"
)

// eventProviderOptionFunc 执行 EventProviderOption 的方法。
type eventProviderOptionFunc func(*eventProviderConfig) *eventProviderConfig

func (o eventProviderOptionFunc) apply(cfg *eventProviderConfig) *eventProviderConfig {
	return o(cfg)
}

// eventProviderConfig EventProvider 初始化配置。
type eventProviderConfig struct {
	Exporters     map[string]EventExporter
	FlushInternal time.Duration
	MaxEvent      int
}

const Internal = 5 * time.Second
const MaxEvent = 9

func newEventProviderConfig(opts ...EventProviderOption) *eventProviderConfig {
	cfg := &eventProviderConfig{
		Exporters:     make(map[string]EventExporter),
		FlushInternal: Internal,
		MaxEvent:      MaxEvent,
	}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// WithExporters 批量设置 EventExporter 。
func WithExporters(exporters ...EventExporter) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		for _, e := range exporters {
			cfg.Exporters[e.Name()] = e
		}
		return cfg
	})
}

// WithServiceInfo 记录服务信息。
func WithServiceInfo(ServiceName string, ServiceVersion string, ServiceInstance string) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		if ServiceName != "" {
			serviceName = ServiceName
		}
		if ServiceVersion != "" {
			serviceVersion = ServiceVersion
		}
		if serviceInstance != "" {
			serviceInstance = ServiceInstance
		}
		return cfg
	})
}

// WithFlushInternal 设置发送间隔。
func WithFlushInternal(flushInternal time.Duration) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		cfg.FlushInternal = flushInternal
		return cfg
	})
}

// WithMaxEvent 设置Event发送上限。
func WithMaxEvent(maxEvent int) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		cfg.MaxEvent = maxEvent
		return cfg
	})
}
