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
	stopCh   chan bool
	stopOnce sync.Once
}

// GetDefaultExporter 获取默认的 EventExporter 。
func GetDefaultExporter() EventExporter {
	return &exporter{
		name:     "DefaultExporter",
		stopCh:   make(chan bool, 1),
		stopOnce: sync.Once{},
	}
}

func (e *exporter) Name() string {
	return e.name
}

func (e *exporter) Shutdown(ctx context.Context) error {
	e.stopOnce.Do(func() {
		e.stopCh <- true
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
	case <-e.stopCh:
		return nil
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
