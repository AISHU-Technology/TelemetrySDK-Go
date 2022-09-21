package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/common"
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

// Shutdown 关闭Exporter，关闭HTTP连接，丢弃导出缓存。
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stopOnce.Do(func() {
		close(e.stopCh)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return e.client.Stop(ctx)
}

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

// NewExporter 创建已启动的Exporter。
func NewExporter(client Client) *Exporter {
	return &Exporter{client: client, stopCh: make(chan struct{})}
}
