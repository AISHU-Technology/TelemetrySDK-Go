package eventsdk

import (
	"context"
	"time"
)

var GlobalEventProvider = NewEventProvider(GetDefaultExporter())

type eventProvider struct {
	Exporter EventExporter
	events   []Event
	Limit    int
	Ticker   *time.Ticker
	sent     chan bool
}

const LIMIT = 3
const TIME = 2 * time.Second

// NewEventProvider 新建 EventProvider 。
func NewEventProvider(exporter EventExporter) EventProvider {
	return &eventProvider{
		Exporter: exporter,
		events:   make([]Event, 0, LIMIT+1),
		Limit:    LIMIT,
		Ticker:   time.NewTicker(TIME),
		sent:     make(chan bool, 10),
	}
}

func GetEventProvider() EventProvider {
	return GlobalEventProvider
}

func SetEventProvider(ep EventProvider) {
	GlobalEventProvider = ep
	go func() {
		ep.SendEvents()
	}()
}

func (ep *eventProvider) LoadEvent(event Event) {
	ep.events = append(ep.events, event)
	ep.verdictSent()
}

func (ep *eventProvider) Shutdown(ctx context.Context) error {
	return ep.Exporter.Shutdown(ctx)
}

func (ep *eventProvider) SendEvents() {
	for {
		select {
		case <-ep.sent:
			ep.ForceFlash()
		case <-ep.Ticker.C:
			ep.ForceFlash()
		}
	}
}

func (ep *eventProvider) ForceFlash() {
	ep.sent = make(chan bool, 10)
	_ = ep.Exporter.ExportEvents(context.Background(), ep.events)
	ep.events = make([]Event, 0, LIMIT+1)
}

func (ep *eventProvider) verdictSent() {
	if len(ep.events) >= ep.Limit {
		ep.sent <- true
	}
}

func (ep *eventProvider) private() {}
