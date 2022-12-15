/*
 * @Author: Nick.nie Nick.nie@aishu.cn
 * @Date: 2022-12-15 03:23:44
 * @LastEditors: Nick.nie Nick.nie@aishu.cn
 * @LastEditTime: 2022-12-15 03:23:45
 * @FilePath: /span/exporter/ar_exporter.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package exporter

import (
	"context"
	"os"
	"sync"
)

// exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type exporter struct {
	name     string
	stopCh   chan struct{}
	stopOnce sync.Once
}

// GetDefaultExporter 获取默认的 EventExporter 。
func GetDefaultExporter() LogExporter {
	return &exporter{
		name:     "DefaultExporter",
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}

func (e *exporter) Name() string {
	return e.name
}

func (e *exporter) Shutdown(ctx context.Context) error {
	// 只关闭一次通道。
	e.stopOnce.Do(func() {
		close(e.stopCh)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (e *exporter) ExportLogs(ctx context.Context, p []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	// 已关闭通道，不发送。
	case <-e.stopCh:
		return nil
	// 正常情况，发送数据。
	default:
		return export(p)
	}

}

// export 执行发送操作，默认发到控制台。
func export(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	//控制台输出
	file := os.Stdout
	_, err := file.Write(p)
	return err
}
