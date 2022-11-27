package eventsdk

import (
	"context"
)

// EventProvider 批量发送数据到 AnyRobot Feed Ingester 的Event数据接收器。
type EventProvider interface {
	// LoadEvent 缓存 Event 等待定时发送或强制发送。
	LoadEvent(event Event)
	// Shutdown 关闭 Event 生产和发送。
	Shutdown(ctx context.Context) error

	SendEvents()

	ForceFlash()

	// private 禁止自己实现接口
	private()
}
