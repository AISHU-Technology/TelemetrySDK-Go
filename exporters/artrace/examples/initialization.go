package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/config"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/errors"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"net/url"
	"time"
)

// Tracer 全局变量用于在业务代码中调用。
var Tracer = otel.GetTracerProvider().Tracer(
	instrumentationName,
	trace.WithInstrumentationVersion(instrumentationVersion),
	trace.WithSchemaURL(instrumentationURL),
)

// AnyRobotURL 从AnyRobot管理端界面获取的Trace上报地址。
var AnyRobotURL = "https://127.0.0.1:6789/traces"

// TracerOption 只能写一个工具库，推荐写使用工具库的根目录。
var (
	instrumentationName    = "go.opentelemetry.io/otel"
	instrumentationVersion = "v1.9.0"
	instrumentationURL     = "https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.9.0"
)

// ServiceResource 用来标记当前服务的信息。
var ServiceResource = resource.NewWithAttributes(
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go/exporters",
	semconv.ServiceNameKey.String("AnyRobotTrace-example"),
	semconv.ServiceVersionKey.String("2.2.0"),
)

// retry 全局变量用于在业务代码中调用。
var retry = config.RetryConfig{
	Enabled:         true,
	InitialInterval: 5 * time.Second,
	MaxInterval:     1 * time.Minute,
	MaxElapsedTime:  5 * time.Minute,
}

// compression 压缩方式：0代表无压缩，1代表GZIP压缩。
var compression int

// InstallExportPipeline 初始化Trace Exporter。
func InstallExportPipeline() (func(context.Context) error, error) {
	if AnyRobotURL == "https://127.0.0.1:6789/traces" {
		return nil, errors.New(customErrors.AnyRobotTraceExporter_UnsetURL)
	}
	Tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(instrumentationURL),
	)
	u, _ := url.Parse(AnyRobotURL)
	c := client.NewHTTPClient(config.WithScheme(u.Scheme), config.WithEndpoint(u.Host),
		config.WithPath(u.Path), config.WithRetry(retry), config.WithCompression(config.Compression(compression)))
	exporter := client.NewExporter(c)

	// Tracer.Start() 启动全局变量Tracer。
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(ServiceResource),
	)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

// SetAnyRobotURL 设置上报地址。
func SetAnyRobotURL(URL string) error {
	_, err := url.Parse(URL)
	if err != nil {
		return errors.New(customErrors.AnyRobotTraceExporter_InvalidURL)
	}
	AnyRobotURL = URL
	return nil
}

// SetInstrumentation 设置调用链主要依赖的工具库。
func SetInstrumentation(InstrumentationName string, InstrumentationVersion string, InstrumentationURL string) error {
	if _, err := url.Parse(InstrumentationURL); err != nil {
		return errors.New(customErrors.AnyRobotTraceExporter_InvalidURL)
	}
	instrumentationName = InstrumentationName
	instrumentationVersion = InstrumentationVersion
	instrumentationURL = InstrumentationURL
	return nil
}

// SetServiceResource 设置当前服务信息。
func SetServiceResource(ServiceResourceURL string, name string, version string) error {
	if _, err := url.Parse(ServiceResourceURL); err != nil {
		return errors.New(customErrors.AnyRobotTraceExporter_InvalidURL)
	}
	ServiceResource = resource.NewWithAttributes(
		ServiceResourceURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)
	return nil
}

// SetRetry 设置重发机制，如果显著干扰到业务运行了，请减少重发时间。
func SetRetry(internal time.Duration, maxInterval time.Duration, maxElapsedTime time.Duration) error {
	if internal > 10*time.Minute || maxInterval > 20*time.Minute || maxElapsedTime > 60*time.Minute {
		return errors.New(customErrors.AnyRobotTraceExporter_RetryTooLong)
	}
	retry = config.RetryConfig{
		Enabled:         true,
		InitialInterval: internal,
		MaxInterval:     maxInterval,
		MaxElapsedTime:  maxElapsedTime,
	}
	return nil
}

// SetCompression 设置压缩方式：0代表无压缩，1代表GZIP压缩。
func SetCompression(Compression int) error {
	if Compression >= 2 || Compression < 0 {
		return errors.New(customErrors.AnyRobotTraceExporter_InvalidCompression)
	}
	compression = Compression
	return nil
}

func NewExporter(c client.Client) *client.Exporter {
	return client.NewExporter(c)
}

func NewStdoutClient() client.Client {
	return client.NewStdoutClient()
}

func NewHTTPClient(opts ...config.HTTPOption) client.Client {
	return client.NewHTTPClient(opts...)
}
