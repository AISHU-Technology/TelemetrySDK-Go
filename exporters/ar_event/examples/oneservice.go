package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

const result = "the answer is"

// addBefore 计算两数之和。
func addBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// add 增加了 Event 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	eventsdk.Info(007, eventsdk.WithEventType("EventExporter/add"))

	//业务代码
	time.Sleep(800 * time.Millisecond)
	return ctx, x + y
}

// multiplyBefore 计算两数之积。
func multiplyBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

// multiply 增加了 Event 的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, span := artrace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(1))
	eventsdk.Warn(800, eventsdk.WithSpanContext(span.SpanContext()))

	//业务代码
	time.Sleep(800 * time.Millisecond)
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
	client := ar_event.NewStdoutClient("")
	exporter := ar_event.NewExporter(client)
	eventProvider := eventsdk.NewEventProvider(eventsdk.WithExporters(exporter))
	eventsdk.SetEventProvider(eventProvider)

	ctx, num := multiply(ctx, 1, 7)
	ctx, num = add(ctx, num, 8)

	defer func() {
		if err := eventProvider.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	ctx := context.Background()
	header := make(map[string]string)
	header["self-defined"] = "some_header"
	client := ar_event.NewHTTPClient(ar_event.WithAnyRobotURL("http://127.0.0.1:8800/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"),
		ar_event.WithCompression(0), ar_event.WithTimeout(10*time.Second), ar_event.WithHeader(header),
		ar_event.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	exporter := ar_event.NewExporter(client)
	eventProvider := eventsdk.NewEventProvider(eventsdk.WithExporters(exporter, eventsdk.GetDefaultExporter()), eventsdk.WithServiceInfo("YourServiceName", "1.0.1", ""))
	eventsdk.SetEventProvider(eventProvider)

	defer func() {
		if err := eventProvider.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
	ctx, num := multiply(ctx, 2, 3)
	for i := 0; i < 12; i++ {
		ctx, num = add(ctx, 2, 3)
	}
	ctx, num = multiply(ctx, 2, 3)
	log.Println(result, num)
}
