package eventsdk

import (
	"context"
	"log"
	"sync"
	"time"
)

// eventProvider 全局唯一，控制 Event 批量发送。
type eventProvider struct {
	RWLock    sync.RWMutex
	Exporters map[string]EventExporter
	Ticker    *time.Ticker
	Limit     int
	Events    []Event
	Sent      chan bool
}

// NewEventProvider 新建 EventProvider 。
func NewEventProvider(opts ...EventProviderOption) EventProvider {
	// 获取默认配置，更新传入配置。
	cfg := defaultProviderConfig
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	// 根据配置创建 EventExporter 。
	return &eventProvider{
		RWLock:    sync.RWMutex{},
		Exporters: cfg.Exporters,
		Ticker:    time.NewTicker(cfg.FlashInternal),
		Limit:     cfg.EventLimit,
		Events:    make([]Event, 0, cfg.EventLimit+1),
		Sent:      make(chan bool, 3),
	}
}

func (ep *eventProvider) Shutdown(ctx context.Context) error {
	var returnErr error = nil
	for _, e := range ep.Exporters {
		err := e.Shutdown(ctx)
		if err != nil {
			if returnErr == nil {
				returnErr = err
			} else {
				log.Println(err)
			}
		}
	}
	return returnErr
}

func (ep *eventProvider) ForceFlash(ctx context.Context) error {
	ep.Sent = make(chan bool, 10)
	if len(ep.Exporters) == 0 {
		ep.Events = make([]Event, 0, ep.Limit+1)
		return nil
	}
	var returnErr error = nil
	for _, e := range ep.Exporters {
		err := e.ExportEvents(ctx, ep.Events)
		if err != nil {
			if returnErr == nil {
				returnErr = err
			} else {
				log.Println(err)
			}
		}
	}

	ep.Events = make([]Event, 0, ep.Limit+1)
	return returnErr
}

func GetEventProvider() EventProvider {
	return globalEventProvider
}

func SetEventProvider(ep EventProvider) {
	globalEventProvider = ep
	go func() {
		ep.(*eventProvider).sendEvents()
	}()
}

// loadEvent 缓存 Event 等待定时发送或强制发送。
func (ep *eventProvider) loadEvent(event Event) {
	ep.Events = append(ep.Events, event)
	ep.verdictSent()
}
func (ep *eventProvider) sendEvents() {
	for {
		select {
		case <-ep.Sent:
			err := ep.ForceFlash(context.Background())
			if err != nil {
				log.Println(err)
			}
		case <-ep.Ticker.C:
			err := ep.ForceFlash(context.Background())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (ep *eventProvider) verdictSent() {
	if len(ep.Events) >= ep.Limit {
		ep.Sent <- true
	}
}

func (ep *eventProvider) private() {}
