package ar_metric

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"encoding/json"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

var _ sdkmetric.Exporter = (*Exporter)(nil)

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

func (e *Exporter) Temporality(k sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.DefaultTemporalitySelector(k)
}
func (e *Exporter) Aggregation(k sdkmetric.InstrumentKind) aggregation.Aggregation {
	return sdkmetric.DefaultAggregationSelector(k)
}
func (e *Exporter) Export(ctx context.Context, data metricdata.ResourceMetrics) error {
	return e.ExportMetrics(ctx, []interface{}{data})
}
func (e *Exporter) ForceFlush(ctx context.Context) error {
	return ctx.Err()
}

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}

// MetricResource 传入 Metric 的默认resource。
func MetricResource() *sdkresource.Resource {
	return resource.MetricResource()
}

var Meter = metric.Meter(nil)
