package ar_event

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/config"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/custom_errors"
	"log"
	"net/url"
	"time"
)

// NewExporter 新建Exporter，需要传入指定的数据发送客户端 client.Client 和名称。
func NewExporter(c client.EventClient) eventsdk.EventExporter {
	return client.NewExporter(c)
}

// NewStdoutClient 创建 client.Exporter 需要的Local数据发送客户端。
func NewStdoutClient(stdoutPath string) client.EventClient {
	return client.NewStdoutClient(stdoutPath)
}

// NewHTTPClient 创建 client.Exporter 需要的HTTP数据发送客户端。
func NewHTTPClient(opts ...config.Option) client.EventClient {
	return client.NewHTTPClient(opts...)
}

// WithAnyRobotURL 设置 client.httpClient 数据上报地址。
func WithAnyRobotURL(URL string) config.Option {
	if _, err := url.Parse(URL); err != nil {
		log.Fatalln(custom_errors.EventExporter_InvalidURL)
		return config.EmptyOption()
	}
	return config.WithAnyRobotURL(URL)
}

// WithCompression 设置 eventsdk.Event 压缩方式：0代表无压缩，1代表GZIP压缩。
func WithCompression(compression int) config.Option {
	if compression >= 2 || compression < 0 {
		log.Fatalln(custom_errors.EventExporter_InvalidCompression)
		return config.EmptyOption()
	}
	return config.WithCompression(config.Compression(compression))
}

// WithTimeout 设置 client.httpClient 连接超时时间。
func WithTimeout(duration time.Duration) config.Option {
	if duration > 60*time.Second || duration < 0 {
		log.Fatalln(custom_errors.EventExporter_DurationTooLong)
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
		log.Fatalln(custom_errors.EventExporter_RetryTooLong)
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
