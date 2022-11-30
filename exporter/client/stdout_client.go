package client

import (
	"context"
	"encoding/json"
	"os"
	"strings"
)

// StdoutClient 客户端结构体。
type StdoutClient struct {
	filepath string
	stopCh   chan struct{}
}

// Path 获取上报地址。
func (d *StdoutClient) Path() string {
	return d.filepath
}

// Stop 关闭发送器。
func (d *StdoutClient) Stop(ctx context.Context) error {
	close(d.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (d *StdoutClient) UploadData(ctx context.Context, data []interface{}) error {
	// 退出逻辑关闭了发送。
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-d.stopCh:
		return nil
	default:

	}
	//控制台输出
	file1 := os.Stdout
	encoder1 := json.NewEncoder(file1)
	encoder1.SetEscapeHTML(false)
	encoder1.SetIndent("", "\t")
	if err := encoder1.Encode(data); err != nil {
		return err
	}
	//写入本地文件，每次覆盖
	if err := encoder1.Encode(data); err != nil {
		return err
	}
	file2, Err := os.Create(d.filepath)
	if Err != nil {
		return Err
	}
	encoder2 := json.NewEncoder(file2)
	encoder2.SetEscapeHTML(false)
	encoder2.SetIndent("", "\t")
	if err := encoder2.Encode(data); err != nil {
		return err
	}
	return nil
}

// NewStdoutClient 创建Exporter的Local客户端。
func NewStdoutClient(stdoutPath string) Client {
	if strings.TrimSpace(stdoutPath) == "" {
		return &StdoutClient{filepath: "./AnyRobotEvent.txt", stopCh: make(chan struct{})}
	}
	return &StdoutClient{filepath: stdoutPath, stopCh: make(chan struct{})}
}
