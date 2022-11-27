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
	FlashInternal time.Duration
	MaxEvent      int
}

func newEventProviderConfig(opts ...EventProviderOption) *eventProviderConfig {
	cfg := &eventProviderConfig{
		Exporters:     make(map[string]EventExporter),
		FlashInternal: 5 * time.Second,
		MaxEvent:      9,
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

// WithFlashInternal 设置发送间隔。
func WithFlashInternal(flashInternal time.Duration) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		cfg.FlashInternal = flashInternal
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
