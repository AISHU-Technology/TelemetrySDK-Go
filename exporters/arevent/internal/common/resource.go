package common

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/version"
	"encoding/json"
	"github.com/shirou/gopsutil/v3/host"
	"net"
)

// Resource 自定义 Event Resource 和 Trace 输出格式一致。
type Resource struct {
	SchemaURL     string                 `json:"SchemaURL"`
	AttributesMap map[string]interface{} `json:"Attributes"`
}

// NewResource 创建新的 *Resource 。
func NewResource() *Resource {
	r := &Resource{
		SchemaURL:     version.SchemaURL,
		AttributesMap: defaultAttributes(),
	}
	return r
}

// MarshalJSON 只输出 AttributesMap 。
func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.AttributesMap)
}

func (r *Resource) GetSchemaURL() string {
	return r.SchemaURL
}

func (r *Resource) GetAttributes() map[string]interface{} {
	return r.AttributesMap
}

// defaultAttributes 获取默认资源信息。
func defaultAttributes() map[string]interface{} {
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
	sdkLanguage := "go"
	sdkMap["language"] = sdkLanguage
	sdkName := "TelemetrySDK-Go/exporters/arevent"
	sdkMap["name"] = sdkName
	sdkVersion := version.VERSION
	sdkMap["version"] = sdkVersion
	// 服务信息
	serviceMap := make(map[string]string, 3)
	result["service"] = serviceMap
	serviceName := "XXXAAA"
	serviceMap["name"] = serviceName
	serviceVersion := "XXXAAA"
	serviceMap["version"] = serviceVersion
	serviceInstance := "XXXAAA"
	serviceMap["instance"] = serviceInstance

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

func (r *Resource) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &r.AttributesMap)
	return err
}
