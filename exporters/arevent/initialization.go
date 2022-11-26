package arevent

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/common"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/errors"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
	"encoding/json"
	"errors"
)

//// Instrumentation 只能记录一个工具库。当前版本不支持修改 Instrumentation 。
//var (
//	instrumentationName    = "TelemetrySDK-Go/exporters/artrace"
//	instrumentationVersion = "v2.2.0"
//	instrumentationURL     = "https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194&path=/exporters/artrace/README.md&_a=preview"
//)

// SetInstrumentation 设置调用链依赖的工具库。
// 当前版本不支持修改Instrumentation。
//func SetInstrumentation(InstrumentationName string, InstrumentationVersion string, InstrumentationURL string) error {
//	if _, err := url.Parse(InstrumentationURL); err != nil {
//		return errors.New(customErrors.AnyRobotEventExporter_InvalidURL)
//	}
//	instrumentationName = InstrumentationName
//	instrumentationVersion = InstrumentationVersion
//	instrumentationURL = InstrumentationURL
//	return nil
//}

//// NewExporter 新建Exporter，需要传入指定的数据发送客户端 client.Client 。
//func NewExporter(c client.Client) *client.Exporter {
//	return client.NewExporter(c)
//}
//
// NewStdoutClient 创建 client.Exporter 需要的Local数据发送客户端。
func NewStdoutClient(stdoutPath string) client.Client {
	return client.NewStdoutClient(stdoutPath)
}

//
//// NewHTTPClient 创建 client.Exporter 需要的HTTP数据发送客户端。
//func NewHTTPClient(opts ...config.Option) client.Client {
//	return client.NewHTTPClient(opts...)
//}
//
//// WithAnyRobotURL 设置 client.httpClient 数据上报地址。
//func WithAnyRobotURL(URL string) config.Option {
//	if _, err := url.Parse(URL); err != nil {
//		log.Fatalln(customErrors.AnyRobotEventExporter_InvalidURL)
//		return config.EmptyOption()
//	}
//	return config.WithAnyRobotURL(URL)
//}
//
//// WithCompression 设置Trace压缩方式：0代表无压缩，1代表GZIP压缩。
//func WithCompression(compression int) config.Option {
//	if compression >= 2 || compression < 0 {
//		log.Fatalln(customErrors.AnyRobotEventExporter_InvalidCompression)
//		return config.EmptyOption()
//	}
//	return config.WithCompression(config.Compression(compression))
//}
//
//// WithTimeout 设置 client.httpClient 连接超时时间。
//func WithTimeout(duration time.Duration) config.Option {
//	if duration > 60*time.Second || duration < 0 {
//		log.Fatalln(customErrors.AnyRobotEventExporter_DurationTooLong)
//		return config.EmptyOption()
//	}
//	return config.WithTimeout(duration)
//}
//
//// WithHeader 设置 client.httpClient 用户自定义请求头。
//func WithHeader(headers map[string]string) config.Option {
//	return config.WithHeader(headers)
//}
//
//// WithRetry 设置 client.httpClient 重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。
//func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, maxElapsedTime time.Duration) config.Option {
//	if enabled && (internal > 10*time.Minute || maxInterval > 20*time.Minute || maxElapsedTime > 60*time.Minute) {
//		log.Fatalln(customErrors.AnyRobotEventExporter_RetryTooLong)
//		return config.EmptyOption()
//	}
//	retry := config.RetryConfig{
//		Enabled:         enabled,
//		InitialInterval: internal,
//		MaxInterval:     maxInterval,
//		MaxElapsedTime:  maxElapsedTime,
//	}
//	return config.WithRetry(retry)
//}

//// GetResource 获取内置资源信息，记录客户服务名，需要传入服务名 serviceName ，服务版本 serviceVersion ，服务实例ID。
//func GetResource(serviceName string, serviceVersion string, serviceInstanceID string) *resource.Resource {
//	//获取主机IP
//	connection, _ := net.Dial("udp", "rockyrori.cn:80")
//	ipPort := connection.LocalAddr().(*net.UDPAddr)
//	hostIP := strings.Split(ipPort.String(), ":")[0]
//	//获取主机信息
//	infoState, _ := host.Info()
//
//	return resource.NewWithAttributes(instrumentationURL)
//}

const (
	ERROR common.Level = common.ERROR
	WARN  common.Level = common.WARN
	INFO  common.Level = common.INFO
)

func NewEvent(eventType string) model.AREvent {
	return common.NewEvent(eventType)
}

// NewAttribute 创建新的 Attribute 。
func NewAttribute(key string, value model.ARValue) model.ARAttribute {
	return common.NewAttribute(key, value)
}

// BoolValue 传入 bool 类型的值。
func BoolValue(value bool) model.ARValue {
	return common.BoolValue(value)
}

// BoolArray 传入 []bool 类型的值。
func BoolArray(value []bool) model.ARValue {
	return common.BoolArray(value)
}

// IntValue 传入 int 类型的值。
func IntValue(value int) model.ARValue {
	return common.IntValue(value)
}

// IntArray 传入 []int 类型的值。
func IntArray(value []int) model.ARValue {
	return common.IntArray(value)
}

// FloatValue 传入 float64 类型的值。
func FloatValue(value float64) model.ARValue {
	return common.FloatValue(value)
}

// FloatArray 传入 []float64 类型的值。
func FloatArray(value []float64) model.ARValue {
	return common.FloatArray(value)
}

// StringValue 传入 string 类型的值。
func StringValue(value string) model.ARValue {
	return common.StringValue(value)
}

// StringArray 传入 []string 类型的值。
func StringArray(value []string) model.ARValue {
	return common.StringArray(value)
}

func UnmarshalEvents(b []byte) ([]model.AREvent, error) {
	events := make([]*common.Event, 0)
	err := json.Unmarshal(b, &events)

	result := make([]model.AREvent, 0, len(events))
	for _, event := range events {
		result = append(result, event)
	}
	if len(result) == 0 {
		err = errors.New(customErrors.AnyRobotEventExporter_InvalidJSON)
	}
	return result, err
}
