package ar_trace

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/common"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"encoding/json"
	"go.opentelemetry.io/otel"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// 跨包实现接口占位用。
var _ sdktrace.SpanExporter = (*TraceExporter)(nil)

// Tracer 是一个全局变量，用于在业务代码中生产Span。
var Tracer = otel.GetTracerProvider().Tracer(
	version.TraceInstrumentationName,
	trace.WithInstrumentationVersion(version.TraceInstrumentationVersion),
	trace.WithSchemaURL(version.TraceInstrumentationURL),
)

// TraceExporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type TraceExporter struct {
	*public.Exporter
}

// ExportSpans 批量发送AnyRobotSpans到AnyRobot Feed Ingester的Trace数据接收器。
func (e *TraceExporter) ExportSpans(ctx context.Context, traces []sdktrace.ReadOnlySpan) error {
	if len(traces) == 0 {
		return nil
	}
	arTrace := common.AnyRobotTraceFromReadOnlyTrace(traces)
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(arTrace); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// NewExporter 创建已启动的 TraceExporter 。
func NewExporter(c public.Client) *TraceExporter {
	return &TraceExporter{
		public.NewExporter(c),
	}
}

// TraceResource 传入 Trace 的默认Resource。
func TraceResource() *sdkresource.Resource {
	return resource.TraceResource()
}
