package eventsdk

import (
	"context"
)

var GlobalEventProvider = NewEventProvider(GetDefaultExporter())

type eventProvider struct {
	Exporter EventExporter
}

// NewEventProvider 新建 EventProvider 。
func NewEventProvider(exporter EventExporter) EventProvider {
	return &eventProvider{
		Exporter: exporter,
	}
}

func GetEventProvider() EventProvider {
	return GlobalEventProvider
}

func SetEventProvider(ep EventProvider) {
	GlobalEventProvider = ep
}

func (p *eventProvider) Load(event Event) {
	events := make([]Event, 0)
	events = append(events, event)
	_ = p.Exporter.ExportEvents(context.Background(), events)
}

func (p *eventProvider) Shutdown(ctx context.Context) error {
	return p.Exporter.Shutdown(ctx)
}
