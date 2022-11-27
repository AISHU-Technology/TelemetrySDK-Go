package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"sync"
)

var _ eventsdk.EventExporter = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type Exporter struct {
	name     string
	client   EventClient
	stopCh   chan struct{}
	stopOnce sync.Once
}

func (e *Exporter) Name() string {
	return e.name
}

// Shutdown 关闭Exporter，关闭HTTP连接，丢弃导出缓存。
func (e *Exporter) Shutdown(ctx context.Context) error {
	var err error = nil
	e.stopOnce.Do(func() {
		close(e.stopCh)
		err = e.client.Stop(ctx)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}

// ExportEvents 批量发送AnyRobotSpans到AnyRobot Feed Ingester的 Event 数据接收器。
func (e *Exporter) ExportEvents(ctx context.Context, events []eventsdk.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.stopCh:
		return nil
	default:
	}
	if len(events) == 0 {
		return nil
	}
	return e.client.UploadEvents(ctx, events)
}

// NewExporter 创建已启动的Exporter。
func NewExporter(client EventClient) *Exporter {
	return &Exporter{
		name:     client.Path(),
		client:   client,
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}
