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
	stopOnce  sync.Once
	stopCh    chan struct{}
}

// NewEventProvider 根据配置项，新建 EventProvider 。
func NewEventProvider(opts ...EventProviderOption) EventProvider {
	// 获取默认配置，更新传入配置。
	cfg := newEventProviderConfig(opts...)
	// 根据配置创建 EventExporter 。
	return &eventProvider{
		RWLock:    sync.RWMutex{},
		Exporters: cfg.Exporters,
		Ticker:    time.NewTicker(cfg.FlashInternal),
		Limit:     cfg.MaxEvent,
		Events:    make([]Event, 0, cfg.MaxEvent+1),
		Sent:      make(chan bool, 3),
		stopOnce:  sync.Once{},
		stopCh:    make(chan struct{}),
	}
}

func (ep *eventProvider) Shutdown(ctx context.Context) error {
	ep.stopOnce.Do(func() {
		close(ep.stopCh)
	})
	// 只返回其中一个错误
	var returnErr error = nil
	for _, e := range ep.Exporters {
		err := e.Shutdown(ctx)
		if err != nil {
			// 如果多余一个错误则记日志
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
	// 调用发送方法，可发送标记位重置。
	ep.Sent = make(chan bool, 3)
	// 不存在发送目标，直接丢弃数据。
	if len(ep.Exporters) == 0 {
		ep.Events = make([]Event, 0, ep.Limit+1)
		return nil
	}
	// 只返回其中一个错误。
	var returnErr error = nil
	// 发送中禁止修改 Events ，上锁。
	ep.RWLock.Lock()
	// 往所有发送地址发送相同的数据。
	for _, e := range ep.Exporters {
		err := e.ExportEvents(ctx, ep.Events)
		if err != nil {
			// 如果多余一个错误则记日志。
			if returnErr == nil {
				returnErr = err
			} else {
				log.Println(err)
			}
		}
	}
	// 发送结束解锁。
	ep.RWLock.Unlock()

	ep.Events = make([]Event, 0, ep.Limit+1)
	return returnErr
}

// GetEventProvider 获取全局唯一 EventProvider 。
func GetEventProvider() EventProvider {
	return globalEventProvider
}

// SetEventProvider 设置全局唯一 EventProvider 。
func SetEventProvider(ep EventProvider) {
	// 先关闭上一个 EventProvider 。
	_ = globalEventProvider.Shutdown(context.Background())
	globalEventProvider = ep
	// 设置了 EventProvider 就开启发送。
	go func() {
		ep.(*eventProvider).sendEvents()
	}()
}

// loadEvent 缓存 Event 等待定时发送或强制发送。
func (ep *eventProvider) loadEvent(event Event) {
	// 增加 Event ，上锁。
	ep.RWLock.Lock()
	ep.Events = append(ep.Events, event)
	// 每次添加 Event ，判断是否超过数量上限。
	ep.verdictSent()
	// 结束载入，解锁。
	ep.RWLock.Unlock()
}

// sendEvents 无限等待定时发送或超过上限发送。
func (ep *eventProvider) sendEvents() {
EXIT:
	for {
		select {
		// 关闭之后退出循环。
		case <-ep.stopCh:
			break EXIT
		// 超过上限发送。
		case <-ep.Sent:
			err := ep.ForceFlash(context.Background())
			if err != nil {
				log.Println(err)
			}
		// 定时发送
		case <-ep.Ticker.C:
			err := ep.ForceFlash(context.Background())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// verdictSent 判断是否超过上限。
func (ep *eventProvider) verdictSent() {
	if len(ep.Events) >= ep.Limit {
		ep.Sent <- true
	}
}

func (ep *eventProvider) private() {}
