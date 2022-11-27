package eventsdk

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

// exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type exporter struct {
	name     string
	stopCh   chan struct{}
	stopOnce sync.Once
}

// GetDefaultExporter 获取默认的 EventExporter 。
func GetDefaultExporter() EventExporter {
	return &exporter{
		name:     "DefaultExporter",
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}

func (e *exporter) Name() string {
	return e.name
}

func (e *exporter) Shutdown(ctx context.Context) error {
	// 只关闭一次通道。
	e.stopOnce.Do(func() {
		close(e.stopCh)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (e *exporter) ExportEvents(ctx context.Context, events []Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// 已关闭通道，不发送。
	case <-e.stopCh:
		return nil
	// 正常情况，发送数据。
	default:
		return export(events)
	}

}

// export 执行发送操作，默认发到控制台。
func export(events []Event) error {
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
