package exporter

import (
	"context"
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
func GetDefaultExporter() LogExporter {
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

func (e *exporter) ExportLogs(ctx context.Context, p []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// 已关闭通道，不发送。
	case <-e.stopCh:
		return nil
	// 正常情况，发送数据。
	default:
		return export(p)
	}

}

// export 执行发送操作，默认发到控制台。
func export(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	//控制台输出
	file := os.Stdout
	_, err := file.Write(p)
	return err
}
