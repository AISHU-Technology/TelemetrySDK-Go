package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

const result = "the answer is"

// addBefore 计算两数之和。
func addBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(300 * time.Millisecond)
	return ctx, x + y
}

// add 增加了Trace埋点的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
	span.SetAttributes(attribute.KeyValue{Key: "add", Value: attribute.StringValue("计算两数之和")})
	span.AddEvent("AddEvent", trace.WithAttributes(attribute.KeyValue{Key: "succeeded", Value: attribute.BoolValue(true)}))
	span.SetStatus(2, "成功计算加法")
	defer span.End()

	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// multiplyBefore 计算两数之积。
func multiplyBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(300 * time.Millisecond)
	return ctx, x * y
}

// multiply 增加了Trace埋点的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(1))
	span.SetAttributes(attribute.KeyValue{Key: "multiply", Value: attribute.StringValue("计算两数之积")})
	span.AddEvent("multiplyEvent", trace.WithAttributes(attribute.String("succeeded", "true"), attribute.String("analyzed", "100ms")))
	span.SetStatus(2, "成功计算乘积")
	defer span.End()

	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

// Example 原始的业务系统入口
func Example() {
	ctx := context.Background()
	ctx, num := multiplyBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)
	ctx, num = addBefore(ctx, num, 8)
	log.Println(result, num)
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
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
	ctx, num = add(ctx, num, 8)
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器输出到Trace接收器。
func HTTPExample() {
	ctx := context.Background()
	client := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/abcd4f634e80d530/events"))
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
	ctx, num = add(ctx, num, 8)
	log.Println(result, num)
}

// HTTPSExample 通过HTTPS发送器输出到Trace接收器。
func HTTPSExample() {
	ctx := context.Background()
	client := artrace.NewHTTPClient(artrace.WithAnyRobotURL("https://a.b.c.d/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"))
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
	ctx, num = add(ctx, num, 8)
	log.Println(result, num)
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	ctx := context.Background()
	header := make(map[string]string)
	header["self-defined"] = "some_header"
	client := artrace.NewHTTPClient(artrace.WithAnyRobotURL("https://a.b.c.d/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"),
		artrace.WithCompression(1), artrace.WithTimeout(10*time.Second), artrace.WithHeader(header),
		artrace.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
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
