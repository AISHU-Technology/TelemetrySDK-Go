# Akashic_TelemetrySDK-Go

[仓库地址](https://devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go?path=%2F&version=GBfeature-arp-205194&_a=contents)

Akashic_TelemetrySDK-Go 是[OpenTelemetry](https://opentelemetry.io/)的[Go](https://golang.org/)
语言版本实现。它提供了一系列接口帮助开发者完成代码埋点过程，旨在提高用户业务的可观测性能力。

## Project Status

| Signal  | Status | Project                                                                                                                                                 |
|---------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| Traces  | Beta   | [trace](https://devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go?version=GBfeature-arp-205194&_a=contents&path=%2Fexporters%2Fartrace) |
| Metrics | Alpha  | N/A                                                                                                                                                     |
| Logs    | Alpha  | [log](https://devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go?version=GBfeature-arp-205194&_a=contents&path=%2Fspan)                  |

## Compatibility

> Akashic_TelemetrySDK-Go 要求Go版本不低于1.17。

## Getting Started

> 引入SDK：go get devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace@feature-arp-205194

> 更新SDK：步骤等同于引入SDK

> 添加依赖：`import (
"context"
"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/anyrobottrace/tracehttp"
"log")`

> 修改入口main.go：

`func main() {
ctx := context.Background()
_ = ar.SetTraceEnvironment()
shutdown, _ := ar.InstallExportPipeline(ctx)
defer func() {
if err := shutdown(ctx); err != nil {
log.Fatal(err)
}
}()
//your code here }`

> 修改业务代码例如logics.go

`func add(ctx context.Context, x, y int64) (context.Context, int64) {
_, span := ar.Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
span.SetAttributes(attribute.KeyValue{Key: "add", Value: attribute.StringValue("计算两数之和")})
defer span.End()
//your code here
return ctx, x + y }`

> 修改settings.go的全部参数，Trace上报地址从前端获取

`func SetTraceEnvironment() {
	_ = setAnyRobotURL("http://10.4.130.68:880/api/feed_ingester/v1/jobs/traceTest/events")
	_ = setInstrumentation("go.opentelemetry.io/otel", "v1.9.0", "https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.9.0")
	_ = setServiceResource("devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go/exporters", "AnyRobotTrace-example", "2.2.0")
	_ = setRetry(5*time.Second, 1*time.Minute, 5*time.Minute)
	_ = setCompression(0)
}`