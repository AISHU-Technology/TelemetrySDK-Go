package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

const result = "the answer is"

// add 增加了Trace埋点的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := ar_trace.Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
	defer span.End()
	span.SetAttributes(attribute.KeyValue{Key: "add", Value: attribute.StringValue("计算两数之和")})
	span.AddEvent("AddEvent", trace.WithAttributes(attribute.KeyValue{Key: "succeeded", Value: attribute.BoolValue(true)}))
	span.SetStatus(2, "成功计算加法")

	// 业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// multiply 增加了Trace埋点的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, span := ar_trace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(1))
	defer span.End()
	span.SetAttributes(attribute.KeyValue{Key: "multiply", Value: attribute.StringValue("计算两数之积")})
	span.AddEvent("multiplyEvent", trace.WithAttributes(attribute.BoolSlice("key", []bool{true, true}), attribute.String("analyzed", "100ms")))
	span.SetStatus(2, "成功计算乘积")

	// 业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

func FileTraceInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	traceClient := public.NewFileClient("./AnyRobotTrace.json")
	traceExporter := ar_trace.NewExporter(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(ar_trace.TraceResource()))
	otel.SetTracerProvider(tracerProvider)
}

func ConsoleTraceInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	traceClient := public.NewConsoleClient()
	traceExporter := ar_trace.NewExporter(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(ar_trace.TraceResource()))
	otel.SetTracerProvider(tracerProvider)
}

func StdoutTraceInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	traceClient := public.NewStdoutClient("./AnyRobotTrace.json")
	traceExporter := ar_trace.NewExporter(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(ar_trace.TraceResource()))
	otel.SetTracerProvider(tracerProvider)
}

func HTTPTraceInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	traceClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-864ab9d78f6a1843/events"))
	traceExporter := ar_trace.NewExporter(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(ar_trace.TraceResource()))
	otel.SetTracerProvider(tracerProvider)
}

func TraceProviderExit(ctx context.Context) {
	tracerProvider := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if err := tracerProvider.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

// FileExample 输出到本地文件。
func FileExample() {
	FileTraceInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	TraceProviderExit(ctx)
	log.Println(result, num)
}

// ConsoleExample 输出到控制台。
func ConsoleExample() {
	ConsoleTraceInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	TraceProviderExit(ctx)
	log.Println(result, num)
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	StdoutTraceInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	TraceProviderExit(ctx)
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器上报到接收器。
func HTTPExample() {
	HTTPTraceInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	TraceProviderExit(ctx)
	log.Println(result, num)
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	ctx := context.Background()
	header := make(map[string]string)
	header["self-defined"] = "some_header"
	traceClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-864ab9d78f6a1843/events"),
		public.WithCompression(1), public.WithTimeout(10*time.Second), public.WithHeader(header),
		public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	traceExporter := ar_trace.NewExporter(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(ar_trace.TraceResource()))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	// 调用ForceFlush之后会立即发送之前生产的2次乘法链路。
	_ = tracerProvider.ForceFlush(ctx)
	// 关闭Trace的发送，这3次加法产生的链路不会发送。
	tracerProvider.UnregisterSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter))
	ctx, num = add(ctx, num, 8)
	ctx, num = add(ctx, num, 9)
	ctx, num = add(ctx, num, 10)
	// 开启Trace的发送，这1次乘法产生的链路会发送。
	tracerProvider.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter))
	ctx, num = multiply(ctx, num, 9)
	log.Println(result, num)
}
