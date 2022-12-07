package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_event"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
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
	eventsdk.Info(struct {
		Name string
		Age  int
	}{"name", 22}, eventsdk.WithEventType("NewEventExporter/add"))

	//业务代码
	time.Sleep(100 * time.Millisecond)
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
	ctx, span := ar_trace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(1))
	eventsdk.Warn(map[string]string{"key": "value", "data": "data"}, eventsdk.WithSpanContext(span.SpanContext()))

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
	eventClient := public.NewStdoutClient("./AnyRobotEvent.txt")
	eventExporter := ar_event.NewEventExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.WithExporters(eventExporter))
	eventsdk.SetEventProvider(eventProvider)

	defer func() {
		if err := eventProvider.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 1, 7)
	ctx, num = add(ctx, num, 8)
	log.Println(result, num)
}

// HTTPExample 修改client所有入参。
func HTTPExample() {
	ctx := context.Background()
	eventClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1:8800/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	eventExporter := ar_event.NewEventExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.WithExporters(eventExporter, eventsdk.GetDefaultExporter()), eventsdk.WithServiceInfo("YourServiceName", "1.0.1", ""))
	eventsdk.SetEventProvider(eventProvider)

	defer func() {
		if err := eventProvider.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	ctx, num := multiply(ctx, 2, 3)
	for i := 0; i < 6; i++ {
		ctx, num = multiply(ctx, 2, 3)
		ctx, num = add(ctx, 2, 3)
	}
	log.Println(result, num)
}
