package eventsdk

import (
	"context"
)

// EventExporter 导出数据到 AnyRobot Feed Ingester 的Event数据接收器。
type EventExporter interface {
	// Shutdown 关闭 EventExporter ，关闭HTTP连接，丢弃缓存数据。
	Shutdown(ctx context.Context) error
	// ExportEvents 批量发送 eventmodel.Event 到 AnyRobot Feed Ingester 的Event数据接收器。
	ExportEvents(ctx context.Context, events []Event) error
}
