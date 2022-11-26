package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
)

// EventClient 负责连接 eventsdk.Event 数据接收器，并且负责转换 eventsdk.Event 数据格式并发送 eventsdk.Event 数据，内部为net/http/client。
type EventClient interface {
	// Stop 用来关闭连接，它只会被调用一次因此不用担心幂等性问题，但是可能存在并发调用，需要上层 EventExporter 通过sync.Once来控制。
	Stop(ctx context.Context) error
	// UploadEvents 用来发送Trace数据，可能会并发调用。
	UploadEvents(ctx context.Context, events []eventsdk.Event) error
}
