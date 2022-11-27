package eventsdk

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

// exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type exporter struct {
	stopCh   chan struct{}
	stopOnce sync.Once
}

func GetDefaultExporter() EventExporter {
	return &exporter{
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}

// ExportEvents 批量发送 AnyRobotEvents 到AnyRobot Feed Ingester的 Event 数据接收器。
func (e *exporter) ExportEvents(ctx context.Context, events []Event) error {
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

	//控制台输出
	file := os.Stdout
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	return encoder.Encode(events)
}

// Shutdown 关闭Exporter，关闭发送，丢弃导出缓存。
func (e *exporter) Shutdown(ctx context.Context) error {
	var err error = nil
	e.stopOnce.Do(func() {
		close(e.stopCh)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return err
}
