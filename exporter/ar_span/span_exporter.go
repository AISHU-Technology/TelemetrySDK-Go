package ar_span

import (
	"context"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
)

var _ encoder.Exporter = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 log 数据接收器。
type Exporter struct {
	publicExporter *public.Exporter
	ctx            context.Context
	cancelFunc     context.CancelFunc
}

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	ctx, cancel := context.WithCancel(context.Background())
	return &Exporter{
		publicExporter: public.NewExporter(c),
		ctx:            ctx,
		cancelFunc:     cancel,
	}
}

// Write 批量发送 log 到AnyRobot Feed Ingester的 log 数据接收器。
func (e *Exporter) Write(p []byte) (n int, err error) {
	exportDataErr := e.publicExporter.ExportData(e.ctx, p)
	if err != nil {
		return 0, exportDataErr
	} else {
		return len(p), nil
	}
}

// 返回cancel函数供sdk结束时调用
func (e *Exporter) GetCancelFunc() context.CancelFunc {
	return e.GetCancelFunc()
}
