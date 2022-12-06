package ar_trace

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace/common"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"encoding/json"
	"github.com/shirou/gopsutil/v3/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"net"
	"strings"
)

var _ sdktrace.SpanExporter = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type Exporter struct {
	*public.Exporter
}

// ExportSpans 批量发送AnyRobotSpans到AnyRobot Feed Ingester的Trace数据接收器。
func (e *Exporter) ExportSpans(ctx context.Context, traces []sdktrace.ReadOnlySpan) error {
	if len(traces) == 0 {
		return nil
	}
	arTrace := common.AnyRobotSpansFromReadOnlySpans(traces)
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(arTrace); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}

// Tracer 是一个全局变量，用于在业务代码中生产Span。
var Tracer = otel.GetTracerProvider().Tracer(
	version.TraceInstrumentationName,
	trace.WithInstrumentationVersion(version.TraceInstrumentationVersion),
	trace.WithSchemaURL(version.TraceInstrumentationURL),
)

// GetResource 获取内置资源信息，记录客户服务名，需要传入服务名 serviceName ，服务版本 serviceVersion ，服务实例ID。
func GetResource(serviceName string, serviceVersion string, serviceInstanceID string) *resource.Resource {
	//获取主机IP
	connection, _ := net.Dial("udp", "255.255.255.255:33")
	ipPort := connection.LocalAddr().(*net.UDPAddr)
	hostIP := strings.Split(ipPort.String(), ":")[0]
	//获取主机信息
	infoState, _ := host.Info()

	return resource.NewWithAttributes(version.TraceInstrumentationURL,
		//主机信息
		semconv.HostNameKey.String(infoState.Hostname),
		semconv.HostArchKey.String(infoState.KernelArch),
		attribute.String("host.ip", hostIP),
		//操作系统信息
		semconv.OSTypeKey.String(infoState.OS),
		semconv.OSDescriptionKey.String(infoState.Platform),
		semconv.OSVersionKey.String(infoState.PlatformVersion),
		//服务信息
		semconv.ServiceInstanceIDKey.String(serviceInstanceID),
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		//版本信息
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKNameKey.String(version.TraceInstrumentationName),
		semconv.TelemetrySDKVersionKey.String(version.TraceInstrumentationVersion),
	)
}
