package exporter

import (
	"context"
	"os"
	"sync"
)

// exporter 导出数据到AnyRobot Feed Ingester的 Log 数据接收器。
type exporter struct {
	name     string
	stopCh   chan struct{}
	stopOnce sync.Once
}

// GetStdoutExporter 获取默认的 LogExporter 。
func GetStdoutExporter() LogExporter {
	return &exporter{
		name:     "StdoutExporter",
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

func (e *exporter) Sync() {
	// 仅用于实现接口，无功能。
}

func SyncStdoutExporter() SyncExporter {
	return &exporter{
		name:     "StdoutExporter",
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}
