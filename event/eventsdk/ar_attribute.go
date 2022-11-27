package eventsdk

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/version"
	"github.com/shirou/gopsutil/v3/host"
	"net"
)

// attribute 自定义 event attribute 和 Trace 输出格式一致。
type attribute struct {
	Key   string `json:"Key"`
	Value value  `json:"Data"`
}

// NewAttribute 创建新的 Attribute 。
func NewAttribute(key string, v Value) Attribute {
	return &attribute{
		Key: key,
		Value: value{
			Type: v.GetType(),
			Data: v.GetData(),
		},
	}
}
func (a *attribute) Valid() bool {
	return a.keyNotNil() && a.keyNotCollide()
}

// keyNotNil 校验 attribute.Key 不为空，即有含义。
func (a *attribute) keyNotNil() bool {
	return len(a.Key) > 0
}

// keyNotCollide 校验 attribute.Key 不与默认值冲突。
func (a *attribute) keyNotCollide() bool {
	switch a.Key {
	case "host":
		return false
	case "os":
		return false
	case "telemetry":
		return false
	case "service":
		return false
	default:
		return true
	}
}

// valueTyped 校验 value.Type 是枚举类型。
//func (a *attribute) valueTyped() bool {
//	switch a.Value.GetType() {
//	case "BOOL":
//		return true
//	case "BOOLARRAY":
//		return true
//	case "INT":
//		return true
//	case "INTARRAY":
//		return true
//	case "FLOAT":
//		return true
//	case "FLOATARRAY":
//		return true
//	case "STRING":
//		return true
//	case "STRINGARRAY":
//		return true
//	default:
//		return false
//	}
//}

func (a *attribute) GetKey() string {
	return a.Key
}

func (a *attribute) GetValue() Value {
	return a.Value
}

func (a *attribute) private() {}

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