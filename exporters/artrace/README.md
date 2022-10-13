# TelemetrySDK-Go

[仓库地址](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194&path=/exporters/artrace/README.md&_a=preview)

`TelemetrySDK-Go`是 [OpenTelemetry](https://opentelemetry.io/) 的 [Go](https://golang.org/)
语言版本实现。本项目提供了一系列接口帮助开发者完成代码埋点过程，旨在提高用户业务的可观测性能力。

## 兼容性

TelemetrySDK-Go 要求Go版本不低于`1.18`。

## 开发指南

1. 检查go版本：`go version`
2. 升级go版本到`go1.18`(https://gomirrors.org/)
3. 引入SDK：
   ```go get devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace@feature-arp-205194```
4. 更新SDK：步骤等同于引入SDK。
5. 新增依赖：以下为新增汇总，以实际使用为准。

```
import (
"context"
"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace"
"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/attribute"
"go.opentelemetry.io/otel/sdk/resource"
sdktrace "go.opentelemetry.io/otel/sdk/trace"
semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
"go.opentelemetry.io/otel/trace"
"log")
```

6. 修改入口：`main.go`

```
func main() {
	ctx := context.Background()
	//c := artrace.NewStdoutClient("")
	c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/abcd4f634e80d530/events"))
	exporter := artrace.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName", "1.0.0")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	//your code here
}
```

7. 修改业务：`logics.go`

```
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "乘法")
	defer span.End()

	//your code here
	return ctx, x * y
}
```

8. 正确填写上报地址：`NewHTTPClient("http://a.b.c.d/")`
   参数从AnyRobot网页端获取。

## 接口文档

// Tracer 全局变量用于在业务代码中调用生产Trace数据。

```
artrace.Tracer.Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span)
```

// GetResource 获取内置资源信息，记录客户服务名。

```
func GetResource(serviceName string, serviceVersion string) *resource.Resource
```

// NewExporter 新建Exporter。

```
func NewExporter(c client.Client) *client.Exporter
```

// NewHTTPClient 创建Exporter的HTTP客户端。

```
func NewHTTPClient(opts ...config.HTTPOption) client.Client
```

// NewStdoutClient 创建Exporter的Local客户端。

```
func NewStdoutClient(stdoutPath string) client.Client
```

// WithAnyRobotURL 设置Trace数据上报地址。

```
func WithAnyRobotURL(URL string) config.HTTPOption
```

// WithCompression 设置压缩方式：0代表无压缩，1代表GZIP压缩。

```
func WithCompression(compression int) config.HTTPOption
```

// WithHeader 设置用户自定义请求头。

```
func WithHeader(headers map[string]string) config.HTTPOption
```

// WithRetry 设置重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。

```
func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, ...) config.HTTPOption
```

// WithTimeout 设置HTTP连接超时时间。

```
func WithTimeout(duration time.Duration) config.HTTPOption
```

## 补充说明

1. 多个服务的改造参考：
   [改造示例](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?path=/exporters/artrace/examples&version=GBfeature-arp-205194)
2. 最佳实践文档：[下载链接暂无]()
3. 