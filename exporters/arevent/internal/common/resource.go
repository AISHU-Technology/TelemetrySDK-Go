package common

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/version"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
	"encoding/json"
	"github.com/shirou/gopsutil/v3/host"
	"net"
)

// Resource 自定义 Event Resource 和 Trace 输出格式一致。
type Resource struct {
	SchemaURL     string                   `json:"SchemaURL"`
	Attributes    []*model.ARAttribute     `json:"Attributes"`
	AttributesMap map[string]model.ARValue `json:"AttributesMap"`
}

// NewResource 创建新的 *Resource 。
func NewResource() *Resource {
	r := &Resource{
		SchemaURL:     "https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?path=/exporters/arevent",
		Attributes:    make([]*model.ARAttribute, 0),
		AttributesMap: defaultAttributes(),
	}
	r.mapToSlice()
	return r
}

//MarshalJSON 只输出 Attributes 。
func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Attributes)
}

// mapToSlice 把 AttributesMap 转成 Attributes 。
func (r *Resource) mapToSlice() {
	result := make([]*model.ARAttribute, 0, len(r.AttributesMap))
	for k, v := range r.AttributesMap {
		s := NewAttribute(k, v)
		result = append(result, &s)

	}
	r.Attributes = result
}

// defaultAttributes 获取默认资源信息。
func defaultAttributes() map[string]model.ARValue {
	ip := getHostIP()
	info := getHostInfo()
	result := make(map[string]model.ARValue)
	// 主机信息
	hostIP := StringValue(ip)
	result["host.ip"] = hostIP
	hostArch := StringValue(info.KernelArch)
	result["host.arch"] = hostArch
	hostName := StringValue(info.Hostname)
	result["host.name"] = hostName
	// 操作系统信息
	osType := StringValue(info.OS)
	result["os.type"] = osType
	osVersion := StringValue(info.PlatformVersion)
	result["os.version"] = osVersion
	osDescription := StringValue(info.Platform)
	result["os.description"] = osDescription
	// 版本信息
	language := StringValue("go")
	result["telemetry.sdk.language"] = language
	name := StringValue("TelemetrySDK-Go/exporters/arevent")
	result["telemetry.sdk.name"] = name
	ver := StringValue(version.VERSION)
	result["telemetry.sdk.version"] = ver

	return result
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

func (r *Resource) GetSchemaURL() string {
	return r.SchemaURL
}

func (r *Resource) GetAttributes() []*model.ARAttribute {
	return r.Attributes
}
