# TelemetrySDK-Go

[仓库地址](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194)

TelemetrySDK-Go 是[OpenTelemetry](https://opentelemetry.io/)的[Go](https://golang.org/)
语言版本实现。它提供了一系列接口帮助开发者完成代码埋点过程，旨在提高用户业务的可观测性能力。

## Project Status

| Signal  | Status | Project                                                                                                                                     |
|---------|--------|---------------------------------------------------------------------------------------------------------------------------------------------|
| Traces  | Beta   | [trace](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194&path=%2Fexporters%2Fartrace) |
| Metrics | Alpha  | N/A                                                                                                                                         |
| Logs    | Alpha  | [log](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GBfeature-arp-205194&path=%2Fspan)                  |

## Compatibility

> TelemetrySDK-Go 要求Go版本不低于1.17。

## Getting Started

> 引入SDK：go get devops.aishu.cn/AISHUDevOps/ONE-Architecture/_
> git/TelemetrySDK-Go.git/exporters/artrace@feature-arp-205194

> 更新SDK：步骤等同于引入SDK

> 添加依赖：
`import (
"context"
"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace"
"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/attribute"
"go.opentelemetry.io/otel/sdk/resource"
sdktrace "go.opentelemetry.io/otel/sdk/trace"
semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
"go.opentelemetry.io/otel/trace"
"log")
> `

> 修改入口main.go：

`
func main() {
ctx := context.Background()
//c := artrace.NewStdoutClient()
c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/traceTest/events"))
exporter := artrace.NewExporter(c)
tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName", "1.0.0")))
otel.SetTracerProvider(tracerProvider)
defer func() {
if err := tracerProvider.Shutdown(ctx); err != nil {
log.Println(err)
}
}()
//your code here }
`

> 修改业务代码例如logics.go

`
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
ctx, span := artrace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(2), trace.WithLinks(trace.Link{}))
span.SetStatus(2, "成功计算乘积")
span.AddEvent("multiplyEvent")
span.SetAttributes(attribute.KeyValue{Key: "succeeded", Value: attribute.BoolValue(true)})
defer span.End()
//your code here
return ctx, x * y
}
`

> 修改NewXXXXClient的入参，Trace上报地址从前端获取

`
c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/traceTest/events"))
`

> 代码埋点可以参考devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/examples