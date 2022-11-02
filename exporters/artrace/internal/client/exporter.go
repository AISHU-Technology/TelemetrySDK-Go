package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/common"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"sync"
)

// Exporter 导出数据到AnyRobot Feed Ingester的Trace数据接收器。
type Exporter struct {
	client   Client
	stopCh   chan struct{}
	stopOnce sync.Once
}

var _ sdktrace.SpanExporter = (*Exporter)(nil)

// ExportSpans 批量发送AnyRobotSpans到AnyRobot Feed Ingester的Trace数据接收器。
func (e *Exporter) ExportSpans(ctx context.Context, ross []sdktrace.ReadOnlySpan) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.stopCh:
		return nil
	default:
	}
	AnyRobotSpans := common.AnyRobotSpansFromReadOnlySpans(ross)
	if len(AnyRobotSpans) == 0 {
		return nil
	}
	return e.client.UploadTraces(ctx, AnyRobotSpans)
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
	}
	return err
}

// NewExporter 创建已启动的Exporter。
func NewExporter(client Client) *Exporter {
	return &Exporter{client: client, stopCh: make(chan struct{})}
}
