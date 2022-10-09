# TelemetrySDK-Go

[仓库地址](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194&path=/exporters/artrace/README.md&_a=preview)

`TelemetrySDK-Go`是[OpenTelemetry](https://opentelemetry.io/)的[Go](https://golang.org/)
语言版本实现。本项目提供了一系列接口帮助开发者完成代码埋点过程，旨在提高用户业务的可观测性能力。

## 兼容性

TelemetrySDK-Go 要求Go版本不低于`1.18`。

## 开发指南

1.引入SDK：
```go get devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace@feature-arp-205194```

1.更新SDK：步骤等同于引入SDK。

2.新增依赖：以下为新增汇总，以实际使用为准。

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

3.修改入口：`main.go`

```
func main() {
	ctx := context.Background()
	//c := artrace.NewStdoutClient()
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

4.修改业务:`logics.go`

```
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(2))
	span.SetStatus(2, "成功计算乘积")
	span.AddEvent("multiplyEvent", trace.WithAttributes(attribute.String("succeeded", "true"), attribute.String("tag2", "something")))
	defer span.End()

	//your code here
	return ctx, x * y
}
```

5.填写上报地址:`NewHTTPClient("http://a.b.c.d/")`

## 接口文档

## 改造示例
不同场景的代码埋点示例可以参考:[devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?path=/exporters/artrace/examples](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?path=/exporters/artrace/examples&version=GBfeature-arp-205194)