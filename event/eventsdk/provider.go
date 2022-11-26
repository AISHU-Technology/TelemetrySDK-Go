package eventsdk

import (
	"context"
)

// EventProvider 批量发送数据到 AnyRobot Feed Ingester 的Event数据接收器。
type EventProvider interface {
	// Load 缓存 eventmodel.Event 等待定时发送或强制发送。
	Load(event Event)
	// Shutdown 关闭 Event 生产和发送。
	Shutdown(ctx context.Context) error
}
