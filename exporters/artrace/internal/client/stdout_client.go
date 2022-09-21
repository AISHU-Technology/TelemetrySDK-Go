package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/common"
	"encoding/json"
	"os"
)

// stdoutClient 客户端结构体。
type stdoutClient struct {
	stopCh chan struct{}
}

// Stop 关闭发送器。
func (d *stdoutClient) Stop(ctx context.Context) error {
	close(d.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// UploadTraces 批量发送Trace数据。
func (d *stdoutClient) UploadTraces(ctx context.Context, AnyRobotSpans []*common.AnyRobotSpan) error {
	ctx.Done()
	//控制台输出
	file1 := os.Stdout
	encoder1 := json.NewEncoder(file1)
	encoder1.SetIndent("", "\t")
	_ = encoder1.Encode(AnyRobotSpans)

	//写入本地文件，每次覆盖
	file2, err := os.Create("./AnyRobotTrace.txt")
	encoder2 := json.NewEncoder(file2)
	encoder2.SetIndent("", "\t")
	_ = encoder2.Encode(AnyRobotSpans)
	return err
}

// NewStdoutClient 创建Exporter的Local客户端。
func NewStdoutClient() Client {
	return &stdoutClient{stopCh: make(chan struct{})}
}
