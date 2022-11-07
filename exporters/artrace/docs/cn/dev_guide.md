# 开发指南

1. 检查go版本：`go version`
2. 升级go版本到：[go1.18](https://gomirrors.org/)
3. 引入SDK：
   ```go get devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace@2.2.0```
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
"go.opentelemetry.io/otel/trace"
"log")
```

6. 修改入口：`main.go`

```
func main() {
	ctx := context.Background()
	//client := artrace.NewStdoutClient("")
	client := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/abcd4f634e80d530/events"))
	exporter := artrace.NewExporter(client)
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

9. (可选)注册/解绑Trace Provider来开启/关闭链路数据的生产和发送：

```
func main() {
	ctx := context.Background()
	client := artrace.NewStdoutClient("./AnyRobotTrace.txt")
	exporter := artrace.NewExporter(client)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName", "1.0.0", "")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	//调用ForceFlush之后会立即发送之前生产的2次乘法链路。
	_ = tracerProvider.ForceFlush(ctx)
	//关闭Trace的发送，这3次加法产生的链路不会发送。
	tracerProvider.UnregisterSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter))
	ctx, num = add(ctx, num, 8)
	ctx, num = add(ctx, num, 9)
	ctx, num = add(ctx, num, 10)
	//开启Trace的发送，这1次乘法产生的链路会发送。
	tracerProvider.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter))
	ctx, num = multiply(ctx, num, 9)
	log.Println(result, num)
}
```
