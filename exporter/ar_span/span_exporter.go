package ar_span

import (
	"context"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
)

var _ exporter.LogExporter = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 log 数据接收器。
type Exporter struct {
	*public.Exporter
}

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}

// ExportLogs 批量发送 log 到AnyRobot Feed Ingester的 log 数据接收器。
func (e *Exporter) ExportLogs(ctx context.Context, p []byte) error {
	return e.ExportData(ctx, p)
}
