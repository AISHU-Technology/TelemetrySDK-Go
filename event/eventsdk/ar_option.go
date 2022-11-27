package eventsdk

import "time"

type eventProviderOptionFunc func(*eventProviderConfig) *eventProviderConfig

func (o eventProviderOptionFunc) apply(cfg *eventProviderConfig) *eventProviderConfig {
	return o(cfg)
}

type eventProviderConfig struct {
	Exporters     map[string]EventExporter
	FlashInternal time.Duration
	EventLimit    int
}

func defaultExporterMap() map[string]EventExporter {
	exporterMap := make(map[string]EventExporter)
	return exporterMap
}

var defaultProviderConfig = &eventProviderConfig{
	Exporters:     defaultExporterMap(),
	FlashInternal: 5 * time.Second,
	EventLimit:    9,
}

func WithExporters(exporters ...EventExporter) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		for _, e := range exporters {
			cfg.Exporters[e.Name()] = e
		}
		return cfg
	})
}
