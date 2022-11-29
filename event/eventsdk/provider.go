package eventsdk

import (
	"context"
)

// EventProvider 批量发送数据到 AnyRobot Feed Ingester 的Event数据接收器。
type EventProvider interface {
	// Shutdown 关闭 Event 生产和发送。
	Shutdown(ctx context.Context) error
	// ForceFlush 立即发送 []Event 。
	ForceFlush(ctx context.Context) error

	// private 禁止自己实现接口
	private()
}
