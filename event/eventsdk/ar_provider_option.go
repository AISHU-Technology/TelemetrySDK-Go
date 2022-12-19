package eventsdk

import (
	"strings"
	"time"
)

// eventProviderOptionFunc 执行 EventProviderOption 的方法。
type eventProviderOptionFunc func(*eventProviderConfig) *eventProviderConfig

func (o eventProviderOptionFunc) apply(cfg *eventProviderConfig) *eventProviderConfig {
	return o(cfg)
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
		if strings.TrimSpace(ServiceName) != "" {
			globalServiceName = ServiceName
		}
		if strings.TrimSpace(globalServiceVersion) != "" {
			globalServiceVersion = ServiceVersion
		}
		if strings.TrimSpace(globalServiceInstance) != "" {
			globalServiceInstance = ServiceInstance
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
