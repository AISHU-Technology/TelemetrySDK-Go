package artrace

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/url"
	"time"
)

// Tracer 全局变量用于在业务代码中调用生产Trace数据。
var Tracer = otel.GetTracerProvider().Tracer(
	instrumentationName,
	trace.WithInstrumentationVersion(instrumentationVersion),
	trace.WithSchemaURL(instrumentationURL),
)

// TracerOption 只能写一个工具库，推荐写使用工具库的根目录。当前版本不支持修改Instrumentation。
var (
	instrumentationName    = "go.opentelemetry.io/otel"
	instrumentationVersion = "v1.9.0"
	instrumentationURL     = "https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.9.0"
)

// SetInstrumentation 设置调用链主要依赖的工具库。
// 当前版本不支持修改Instrumentation。
//func SetInstrumentation(InstrumentationName string, InstrumentationVersion string, InstrumentationURL string) error {
//	if _, err := url.Parse(InstrumentationURL); err != nil {
//		return errors.New(customErrors.AnyRobotTraceExporter_InvalidURL)
//	}
//	instrumentationName = InstrumentationName
//	instrumentationVersion = InstrumentationVersion
//	instrumentationURL = InstrumentationURL
//	return nil
//}

// NewExporter 创建已启动的Exporter。
func NewExporter(c client.Client) *client.Exporter {
	return client.NewExporter(c)
}

// NewStdoutClient 创建Exporter的Local客户端。
func NewStdoutClient(stdoutPath string) client.Client {
	return client.NewStdoutClient(stdoutPath)
}

// NewHTTPClient 创建Exporter的HTTP客户端。
func NewHTTPClient(opts ...config.HTTPOption) client.Client {
	return client.NewHTTPClient(opts...)
}

// WithAnyRobotURL 设置Trace数据上报地址。
func WithAnyRobotURL(URL string) config.HTTPOption {
	if _, err := url.Parse(URL); err != nil {
		log.Println(customErrors.AnyRobotTraceExporter_InvalidURL)
	}
	return config.WithAnyRobotURL(URL)
}

// WithCompression 设置压缩方式：0代表无压缩，1代表GZIP压缩。
func WithCompression(compression int) config.HTTPOption {
	if compression >= 2 || compression < 0 {
		log.Println(customErrors.AnyRobotTraceExporter_InvalidCompression)
	}
	return config.WithCompression(config.Compression(compression))
}

// WithTimeout 设置HTTP连接超时时间。
func WithTimeout(duration time.Duration) config.HTTPOption {
	return config.WithTimeout(duration)
}

// WithHeader 设置用户自定义请求头。
func WithHeader(headers map[string]string) config.HTTPOption {
	return config.WithHeader(headers)
}

// WithRetry 设置重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。
func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, maxElapsedTime time.Duration) config.HTTPOption {
	if internal > 10*time.Minute || maxInterval > 20*time.Minute || maxElapsedTime > 60*time.Minute {
		log.Println(customErrors.AnyRobotTraceExporter_RetryTooLong)
	}
	retry := config.RetryConfig{
		Enabled:         enabled,
		InitialInterval: internal,
		MaxInterval:     maxInterval,
		MaxElapsedTime:  maxElapsedTime,
	}
	return config.WithRetry(retry)
}

// GetResource 获取内置资源信息，记录客户服务名。
func GetResource(serviceName string, serviceVersion string) *resource.Resource {
	return resource.NewWithAttributes("devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace",
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
		attribute.String("telemetry.sdk.language", "go"),
		attribute.String("telemetry.sdk.name", "ONE-Architecture"),
		attribute.String("telemetry.sdk.version", "2.2.0"),
	)
}
