package examples

import (
	"context"
	"crypto/tls"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/url"
	"strings"
	"time"
)

const result = "the answer is"

// add 计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
	span.SetAttributes(attribute.KeyValue{Key: "add", Value: attribute.StringValue("计算两数之和")})
	defer span.End()

	time.Sleep(300 * time.Millisecond)

	return ctx, x + y
}

// multiply 计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := Tracer.Start(ctx, "乘法", trace.WithSpanKind(2), trace.WithLinks(trace.Link{}))
	span.SetStatus(2, "成功计算乘积")
	span.AddEvent("multiplyEvent")
	span.SetAttributes(attribute.KeyValue{Key: "succeeded", Value: attribute.BoolValue(true)})
	span.SetAttributes(attribute.KeyValue{Key: "succeeded", Value: attribute.BoolSliceValue([]bool{})})
	defer span.End()

	time.Sleep(300 * time.Millisecond)
	return ctx, x * y
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	ctx := context.Background()
	c := NewStdoutClient()
	exporter := NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(ServiceResource))
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
	u, _ := url.Parse(strings.TrimSpace("http://10.4.130.68:880/api/feed_ingester/v1/jobs/traceTest/events"))
	c := NewHTTPClient(config.WithScheme(u.Scheme), config.WithEndpoint(u.Host),
		config.WithPath(u.Path), config.WithCompression(config.GzipCompression))
	exporter := client.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(ServiceResource))
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
	u, _ := url.Parse(strings.TrimSpace("https://10.4.107.107/api/feed_ingester/v1/jobs/job-a6d44f634e80d530/events"))
	c := NewHTTPClient(config.WithScheme(u.Scheme), config.WithEndpoint(u.Host),
		config.WithPath(u.Path), config.WithCompression(config.GzipCompression), config.WithTimeout(10*time.Second),
		config.WithHeader(nil), config.WithTLSClientConfig(&tls.Config{InsecureSkipVerify: true}))
	exporter := client.NewExporter(c)
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(ServiceResource))
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
