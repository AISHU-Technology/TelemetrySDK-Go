package public

import (
	"context"
	"os"
	"strings"
)

// StdoutClient 客户端结构体。
type StdoutClient struct {
	filepath string
	stopCh   chan struct{}
}

// Path 获取上报地址。
func (c *StdoutClient) Path() string {
	return c.filepath
}

// Stop 关闭发送器。
func (c *StdoutClient) Stop(ctx context.Context) error {
	close(c.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (c *StdoutClient) UploadData(ctx context.Context, data []byte) error {
	// 退出逻辑关闭了发送。
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.stopCh:
		return nil
	default:

	}
	//控制台输出
	file1 := os.Stdout
	if _, err := file1.Write(data); err != nil {
		return err
	}
	//写入本地文件，每次覆盖
	file2, Err := os.Create(c.filepath)
	if Err != nil {
		return Err
	}
	if _, err := file2.Write(data); err != nil {
		return err
	}
	return nil
}

// NewStdoutClient 创建Exporter的Local客户端。
func NewStdoutClient(stdoutPath string) Client {
	if strings.TrimSpace(stdoutPath) == "" {
		return &StdoutClient{filepath: "./AnyRobotData.txt", stopCh: make(chan struct{})}
	}
	return &StdoutClient{filepath: stdoutPath, stopCh: make(chan struct{})}
}
