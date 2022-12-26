package resource

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"github.com/shirou/gopsutil/v3/host"
	"go.opentelemetry.io/otel/attribute"
	environment "go.opentelemetry.io/otel/sdk/resource"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"net"
	"strings"
)

var (
	globalServiceName     = defaultServiceName()
	globalServiceVersion  = "UnknownServiceVersion"
	globalServiceInstance = "UnknownServiceInstance"
)

func SetServiceName(serviceName string) {
	globalServiceName = serviceName
}

func SetServiceVersion(serviceVersion string) {
	globalServiceVersion = serviceVersion
}

func SetServiceInstance(serviceInstance string) {
	globalServiceInstance = serviceInstance
}

func GetServiceName() string {
	return globalServiceName
}

func GetServiceVersion() string {
	return globalServiceVersion
}

func GetServiceInstance() string {
	return globalServiceInstance
}

// 当服务名未设置时，自动获取一个默认的服务名。
func defaultServiceName() string {
	attributes := environment.Default().Attributes()
	if len(attributes) > 0 {
		if attributes[0].Key == "service.name" {
			if v := strings.Split(attributes[0].Value.AsString(), "___"); len(v) >= 2 {
				return strings.Split(attributes[0].Value.AsString(), "___")[1]
			}
		}
	}
	return "UnknownServiceName"
}

// getHostIP 获取主机IP。
func getHostIP() string {
	connection, _ := net.Dial("udp", "255.255.255.255:33")
	ipPort := connection.LocalAddr().(*net.UDPAddr)
	return ipPort.IP.String()
}

// getHostInfo 获取主机信息。
func getHostInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}

// innerAttributes 记录Resource中的公共部分。
func innerAttributes() []attribute.KeyValue {
	// 获取本机IP。
	ip := getHostIP()
	info := getHostInfo()
	var inner = make([]attribute.KeyValue, 0, 9)
	// 主机信息。
	inner = append(inner, attribute.String("host.ip", ip))
	inner = append(inner, semconv.HostArchKey.String(info.KernelArch))
	inner = append(inner, semconv.HostNameKey.String(info.Hostname))
	// 操作系统信息。
	inner = append(inner, semconv.OSTypeKey.String(info.OS))
	inner = append(inner, semconv.OSVersionKey.String(info.PlatformVersion))
	inner = append(inner, semconv.OSDescriptionKey.String(info.Platform))
	// 服务信息。
	inner = append(inner, semconv.ServiceNameKey.String(GetServiceName()))
	inner = append(inner, semconv.ServiceVersionKey.String(GetServiceVersion()))
	inner = append(inner, semconv.ServiceInstanceIDKey.String(GetServiceInstance()))
	return inner
}

func TraceResource() *sdkresource.Resource {
	var attributes = innerAttributes()
	attributes = append(attributes, semconv.TelemetrySDKLanguageGo)
	attributes = append(attributes, semconv.TelemetrySDKNameKey.String(version.TraceInstrumentationName))
	attributes = append(attributes, semconv.TelemetrySDKVersionKey.String(version.TraceInstrumentationVersion))
	return sdkresource.NewWithAttributes(version.TraceInstrumentationURL, attributes...)
}

func MetricResource() *sdkresource.Resource {
	var attributes = innerAttributes()
	attributes = append(attributes, semconv.TelemetrySDKLanguageGo)
	attributes = append(attributes, semconv.TelemetrySDKNameKey.String(version.MetricInstrumentationName))
	attributes = append(attributes, semconv.TelemetrySDKVersionKey.String(version.MetricInstrumentationVersion))
	return sdkresource.NewWithAttributes(version.MetricInstrumentationURL, attributes...)
}

func EventResource() eventsdk.EventProviderOption {
	return eventsdk.WithServiceInfo(GetServiceName(), GetServiceVersion(), GetServiceInstance())
}