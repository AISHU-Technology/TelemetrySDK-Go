package ar_span

import (
	"context"
	spanLog "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/log"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
)

// 跨包实现接口占位用。
var _ exporter.LogExporter = (*SpanExporter)(nil)

// system_logger 程序日志记录器。
var system_logger *spanLog.SamplerLogger

// service_logger 业务日志记录器。
var service_logger *spanLog.SamplerLogger

// SpanExporter 导出数据到AnyRobot Feed Ingester的 Log 数据接收器。
type SpanExporter struct {
	*public.Exporter
}

// ExportLogs 批量发送 log 到AnyRobot Feed Ingester的 Log 数据接收器。
func (e *SpanExporter) ExportLogs(ctx context.Context, logs []byte) error {
	return e.ExportData(ctx, logs)
}

// NewExporter 创建已启动的 SpanExporter。
func NewExporter(c public.Client) *SpanExporter {
	return &SpanExporter{
		public.NewExporter(c),
	}
}
