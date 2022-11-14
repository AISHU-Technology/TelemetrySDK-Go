package artrace

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/errors"
	"github.com/shirou/gopsutil/v3/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

// Tracer 是一个全局变量，用于在业务代码中生产Span。
var Tracer = otel.GetTracerProvider().Tracer(
	instrumentationName,
	trace.WithInstrumentationVersion(instrumentationVersion),
	trace.WithSchemaURL(instrumentationURL),
)

// Instrumentation 只能记录一个工具库。当前版本不支持修改 Instrumentation 。
var (
	instrumentationName    = "TelemetrySDK-Go/exporters/artrace"
	instrumentationVersion = "v2.2.0"
	instrumentationURL     = "https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?path=/exporters/artrace"
)

// SetInstrumentation 设置调用链依赖的工具库。
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

// NewExporter 新建Exporter，需要传入指定的数据发送客户端 client.Client 。
func NewExporter(c client.Client) *client.Exporter {
	return client.NewExporter(c)
}

// NewStdoutClient 创建 client.Exporter 需要的Local数据发送客户端。
func NewStdoutClient(stdoutPath string) client.Client {
	return client.NewStdoutClient(stdoutPath)
}

// NewHTTPClient 创建 client.Exporter 需要的HTTP数据发送客户端。
func NewHTTPClient(opts ...config.Option) client.Client {
	return client.NewHTTPClient(opts...)
}

// WithAnyRobotURL 设置 client.httpClient 数据上报地址。
func WithAnyRobotURL(URL string) config.Option {
	if _, err := url.Parse(URL); err != nil {
		log.Fatalln(customErrors.AnyRobotTraceExporter_InvalidURL)
		return config.EmptyOption()
	}
	return config.WithAnyRobotURL(URL)
}

// WithCompression 设置Trace压缩方式：0代表无压缩，1代表GZIP压缩。
func WithCompression(compression int) config.Option {
	if compression >= 2 || compression < 0 {
		log.Fatalln(customErrors.AnyRobotTraceExporter_InvalidCompression)
		return config.EmptyOption()
	}
	return config.WithCompression(config.Compression(compression))
}

// WithTimeout 设置 client.httpClient 连接超时时间。
func WithTimeout(duration time.Duration) config.Option {
	if duration > 60*time.Second || duration < 0 {
		log.Fatalln(customErrors.AnyRobotTraceExporter_DurationTooLong)
		return config.EmptyOption()
	}
	return config.WithTimeout(duration)
}

// WithHeader 设置 client.httpClient 用户自定义请求头。
func WithHeader(headers map[string]string) config.Option {
	return config.WithHeader(headers)
}

// WithRetry 设置 client.httpClient 重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。
func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, maxElapsedTime time.Duration) config.Option {
	if enabled && (internal > 10*time.Minute || maxInterval > 20*time.Minute || maxElapsedTime > 60*time.Minute) {
		log.Fatalln(customErrors.AnyRobotTraceExporter_RetryTooLong)
		return config.EmptyOption()
	}
	retry := config.RetryConfig{
		Enabled:         enabled,
		InitialInterval: internal,
		MaxInterval:     maxInterval,
		MaxElapsedTime:  maxElapsedTime,
	}
	return config.WithRetry(retry)
}

// GetResource 获取内置资源信息，记录客户服务名，需要传入服务名 serviceName ，服务版本 serviceVersion ，服务实例ID。
func GetResource(serviceName string, serviceVersion string, serviceInstanceID string) *resource.Resource {
	//获取主机IP
	connection, _ := net.Dial("udp", "255.255.255.255:33")
	ipPort := connection.LocalAddr().(*net.UDPAddr)
	hostIP := strings.Split(ipPort.String(), ":")[0]
	//获取主机信息
	infoState, _ := host.Info()

	return resource.NewWithAttributes(instrumentationURL,
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
		semconv.TelemetrySDKNameKey.String(instrumentationName),
		semconv.TelemetrySDKVersionKey.String(instrumentationVersion),
	)
}
