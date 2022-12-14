/*
 * @Author: Nick.nie Nick.nie@aishu.cn
 * @Date: 2022-12-09 03:07:50
 * @LastEditors: Nick.nie Nick.nie@aishu.cn
 * @LastEditTime: 2022-12-14 02:56:24
 * @FilePath: /span/open_standard/opentelemetry.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package open_standard

import (
	"net"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/config"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	"github.com/shirou/gopsutil/v3/host"
)

const (
	rootSpan = iota
	//OpenTelemetrySDKVersion = "v1.6.1"
	SDKName     = "TelemetrySDK-Go/span"
	SDKVersion  = "2.0.1"
	SDKLanguage = "go"

	serviceName     = "UnknownServiceName"
	serviceVersion  = "UnknownServiceVersion"
	serviceInstance = "UnknownServiceInstance"
)

type Writer interface {
	Write([]field.LogSpan) error
	Close() error
}

type OpenTelemetry struct {
	Encoder  encoder.Encoder
	Resource field.Field
}

func NewOpenTelemetry(enc encoder.Encoder, resources field.Field) OpenTelemetry {
	res := OpenTelemetry{
		Encoder:  enc,
		Resource: resources,
	}
	if res.Resource == nil {
		res.SetDefaultResources()
	}

	return res
}

func (o *OpenTelemetry) Write(t []field.LogSpan) error {
	return o.write(t, rootSpan)
}

func (o *OpenTelemetry) SetDefaultResources() {
	defaultResource := getDefaultResource()
	service := make(map[string]interface{})
	service["name"] = serviceName
	service["version"] = serviceVersion
	service["instance"] = map[string]string{"id": serviceInstance}
	defaultResource["service"] = service
	o.Resource = field.MapField(defaultResource)
}

func (o *OpenTelemetry) SetResourcesWithServiceInfo(ServiceName string, ServiceVersion string, ServiceInstanceID string) {
	defaultResource := getDefaultResource()
	service := make(map[string]interface{})
	service["name"] = ServiceName
	service["version"] = ServiceVersion
	service["instance"] = map[string]string{"id": ServiceInstanceID}
	defaultResource["service"] = service
	o.Resource = field.MapField(defaultResource)
}

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

// getDefaultAttributes 获取默认资源信息。
func getDefaultResource() map[string]interface{} {
	// 获取本机IP
	ip := getHostIP()
	info := getHostInfo()
	result := make(map[string]interface{})
	// 主机信息
	hostMap := make(map[string]string, 3)
	result["host"] = hostMap
	hostIP := ip
	hostMap["ip"] = hostIP
	hostArch := info.KernelArch
	hostMap["arch"] = hostArch
	hostName := info.Hostname
	hostMap["name"] = hostName
	// 操作系统信息
	osMap := make(map[string]string, 3)
	result["os"] = osMap
	osType := info.OS
	osMap["type"] = osType
	osVersion := info.PlatformVersion
	osMap["version"] = osVersion
	osDescription := info.Platform
	osMap["description"] = osDescription
	// 版本信息
	sdkMap := make(map[string]string, 3)
	telemetryMap := make(map[string]interface{}, 1)
	telemetryMap["sdk"] = sdkMap
	result["telemetry"] = telemetryMap
	sdkMap["language"] = SDKLanguage
	sdkMap["name"] = SDKName
	sdkMap["version"] = SDKVersion
	return result
}

func (o *OpenTelemetry) Close() error {
	return o.Encoder.Close()
}

func (o *OpenTelemetry) write(logSpans []field.LogSpan, flag int) error {
	var err error
	telemetrys := field.MallocArrayField(config.MaxLog + 1)
	for _, t := range logSpans {
		telemetry := field.MallocStructField(8)

		link := field.MallocStructField(2)
		link.Set("TraceId", field.StringField(t.TraceID()))
		link.Set("SpanId", field.StringField(t.SpanID()))

		telemetry.Set("Link", link)
		telemetry.Set("Timestamp", field.StringField(time.Now().Format(time.RFC3339Nano)))
		telemetry.Set("SeverityText", t.GetLogLevel())

		telemetry.Set("Body", t.GetRecord())
		attrs := t.GetAttributes()

		telemetry.Set("Attributes", attrs)

		if o.Resource == nil {
			o.SetDefaultResources()
		}
		o.dealResource()
		telemetry.Set("Resource", o.Resource)
		telemetrys.Append(telemetry)
	}

	err = o.Encoder.Write(telemetrys)
	if err != nil {
		return err
	}

	return err
}

func (o *OpenTelemetry) dealResource() {
	resMap, ok := o.Resource.(field.MapField)
	if ok {
		_, serviceInfoOk := resMap["service"]
		_, telemetryInfoOk := resMap["telemetry"]
		_, hostInfoOk := resMap["host"]
		_, osInfoOk := resMap["os"]
		if serviceInfoOk {
			if telemetryInfoOk && hostInfoOk && osInfoOk {
				return
			} else {
				defaultResource := getDefaultResource()
				for k, v := range defaultResource {
					resMap[k] = v
				}
				o.Resource = field.MapField(resMap)
				return
			}
		}

	}
	defaultResource := getDefaultResource()
	service := make(map[string]interface{})
	service["name"] = serviceName
	service["version"] = serviceVersion
	service["instance"] = map[string]string{"id": serviceInstance}
	defaultResource["service"] = service
	defaultResource["customer"] = o.Resource
	o.Resource = field.MapField(defaultResource)
}
