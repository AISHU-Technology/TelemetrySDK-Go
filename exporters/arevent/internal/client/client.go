package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
)

// Client 负责连接Trace数据接收器，并且负责转换Trace数据格式并发送Trace数据，内部为net/http/client。
type Client interface {
	// Stop 用来关闭连接，它只会被调用一次因此不用担心幂等性问题，但是可能存在并发调用，需要上层Exporter通过sync.Once来控制。
	Stop(ctx context.Context) error
	// UploadTraces 用来发送Trace数据，可能会并发调用。
	UploadEvents(ctx context.Context, AnyRobotEvents []model.AREvent) error

	UploadEvent(ctx context.Context, AnyRobotEvent model.AREvent) error
}
