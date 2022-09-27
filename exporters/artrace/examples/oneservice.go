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

// add 计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
	span.SetAttributes(attribute.KeyValue{Key: "add", Value: attribute.StringValue("计算两数之和")})
	span.SetStatus(2, "成功计算两数之和")
	span.AddEvent("addEvent", trace.WithAttributes(attribute.KeyValue{
		Key:   "succeeded",
		Value: attribute.BoolValue(true),
	}))
	defer span.End()

	time.Sleep(300 * time.Millisecond)

	return ctx, x + y
}

// multiply 计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := artrace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(2))
	span.SetStatus(2, "成功计算乘积")
	span.AddEvent("multiplyEvent", trace.WithAttributes(attribute.String("succeeded", "true"), attribute.String("tag2", "something")))
	defer span.End()

	time.Sleep(300 * time.Millisecond)
	return ctx, x * y
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	ctx := context.Background()
	c := artrace.NewStdoutClient()
	exporter := artrace.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 2)
	ctx, num = multiply(ctx, num, 10)
	ctx, num = add(ctx, num, 2)
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器输出到Trace接收器。
func HTTPExample() {
	ctx := context.Background()
	c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("http://a.b.c.d/api/feed_ingester/v1/jobs/abcd4f634e80d530/events"))
	exporter := artrace.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 2)
	ctx, num = multiply(ctx, num, 10)
	ctx, num = add(ctx, num, 2)
	log.Println(result, num)
}

// HTTPSExample 通过HTTPS发送器输出到Trace接收器。
func HTTPSExample() {
	ctx := context.Background()
	c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("https://a.b.c.d/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"))
	exporter := artrace.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 2)
	ctx, num = multiply(ctx, num, 10)
	ctx, num = add(ctx, num, 2)
	log.Println(result, num)
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	ctx := context.Background()
	header := make(map[string]string)
	header["self-defined"] = "some_header"
	c := artrace.NewHTTPClient(artrace.WithAnyRobotURL("https://a.b.c.d/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"),
		artrace.WithCompression(1), artrace.WithTimeout(10*time.Second), artrace.WithHeader(header),
		artrace.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	exporter := artrace.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(artrace.GetResource("YourServiceName")))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 2)
	ctx, num = multiply(ctx, num, 10)
	ctx, num = add(ctx, num, 2)
	log.Println(result, num)
}
