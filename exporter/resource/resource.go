package resource

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	"github.com/shirou/gopsutil/v3/host"
	"go.opentelemetry.io/otel/attribute"
	environment "go.opentelemetry.io/otel/sdk/resource"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"net"
	"strings"
)

// 用户未设置服务信息时的默认值，服务名可以自动获取。
var (
	globalServiceName     = defaultServiceName()
	globalServiceVersion  = "UnknownServiceVersion"
	globalServiceInstance = "UnknownServiceInstance"
)

// SetServiceName 设置服务名
func SetServiceName(serviceName string) {
	globalServiceName = serviceName
}

// SetServiceVersion 设置服务版本
func SetServiceVersion(serviceVersion string) {
	globalServiceVersion = serviceVersion
}

// SetServiceInstance 设置服务实例ID
func SetServiceInstance(serviceInstance string) {
	globalServiceInstance = serviceInstance
}

// GetServiceName 获取服务名
func GetServiceName() string {
	return globalServiceName
}

// GetServiceVersion 获取服务版本
func GetServiceVersion() string {
	return globalServiceVersion
}

// GetServiceInstance 获取服务实例ID
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

// TraceResource 填充Trace资源信息
func TraceResource() *sdkresource.Resource {
	var attributes = innerAttributes()
	attributes = append(attributes, semconv.TelemetrySDKLanguageGo)
	attributes = append(attributes, semconv.TelemetrySDKNameKey.String(version.TraceInstrumentationName))
	attributes = append(attributes, semconv.TelemetrySDKVersionKey.String(version.TraceInstrumentationVersion))
	return sdkresource.NewWithAttributes(version.TraceInstrumentationURL, attributes...)
}

// LogResource 填充Log资源信息
func LogResource() field.Field {
	return field.WithServiceInfo(GetServiceName(), GetServiceVersion(), GetServiceInstance())
}

// MetricResource 填充Metric资源信息
func MetricResource() *sdkresource.Resource {
	var attributes = innerAttributes()
	attributes = append(attributes, semconv.TelemetrySDKLanguageGo)
	attributes = append(attributes, semconv.TelemetrySDKNameKey.String(version.MetricInstrumentationName))
	attributes = append(attributes, semconv.TelemetrySDKVersionKey.String(version.MetricInstrumentationVersion))
	return sdkresource.NewWithAttributes(version.MetricInstrumentationURL, attributes...)
}

// EventResource 填充Event资源信息
func EventResource() eventsdk.EventProviderOption {
	return eventsdk.ServiceInfo(GetServiceName(), GetServiceVersion(), GetServiceInstance())
}
