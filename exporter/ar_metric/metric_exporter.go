package ar_metric

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"encoding/json"
)

// Exporter 导出数据到AnyRobot Feed Ingester的 Metric 数据接收器。
type Exporter struct {
	*public.Exporter
}

// ExportMetrics 批量发送 AnyRobotMetrics 到AnyRobot Feed Ingester的 Metric 数据接收器。
func (e *Exporter) ExportMetrics(ctx context.Context, metrics []interface{}) error {
	if len(metrics) == 0 {
		return nil
	}
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(metrics); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// NewMetricExporter 创建已启动的Exporter。
func NewMetricExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}
