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

// add 增加了 Event 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	eventsdk.Info(struct {
		Name string
		Age  int
	}{"name", 31},
		eventsdk.WithEventType("NewExporter/add"),
	)

	// 业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// multiply 增加了 Event 的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, span := ar_trace.Tracer.Start(ctx, "乘法", trace.WithSpanKind(1))
	eventsdk.Warn(map[string]string{"key": "value", "data": "data"},
		eventsdk.WithSpanContext(span.SpanContext()),
		eventsdk.WithAttributes(eventsdk.NewAttribute("key", false)),
		eventsdk.WithSubject("主题"),
	)

	// 业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

func FileEventInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	eventClient := public.NewFileClient("./AnyRobotEvent.json")
	eventExporter := ar_event.NewExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.Exporters(eventExporter), ar_event.EventResource())
	eventsdk.SetEventProvider(eventProvider)
}

func ConsoleEventInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	eventClient := public.NewConsoleClient()
	eventExporter := ar_event.NewExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.Exporters(eventExporter), ar_event.EventResource())
	eventsdk.SetEventProvider(eventProvider)
}

func StdoutEventInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	eventClient := public.NewStdoutClient("./AnyRobotEvent.json")
	eventExporter := ar_event.NewExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.Exporters(eventExporter), ar_event.EventResource())
	eventsdk.SetEventProvider(eventProvider)
}

func HTTPEventInit() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	eventClient := public.NewHTTPClient(
		public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-864ab9d78f6a1843/events"),
		public.WithCompression(0),
		public.WithTimeout(10*time.Second),
		public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute),
	)
	eventExporter := ar_event.NewExporter(eventClient)
	eventProvider := eventsdk.NewEventProvider(eventsdk.Exporters(eventExporter, eventsdk.GetDefaultExporter()), ar_event.EventResource())
	eventsdk.SetEventProvider(eventProvider)
}

func EventProviderExit() {
	if err := eventsdk.GetEventProvider().Shutdown(); err != nil {
		log.Println(err)
	}
}

// FileExample 输出到本地文件。
func FileExample() {
	FileEventInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	EventProviderExit()
	log.Println(result, num)
}

// ConsoleExample 输出到控制台。
func ConsoleExample() {
	ConsoleEventInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	EventProviderExit()
	log.Println(result, num)
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	StdoutEventInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	EventProviderExit()
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器上报到接收器。
func HTTPExample() {
	HTTPEventInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	EventProviderExit()
	log.Println(result, num)
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	public.SetServiceInfo("YourServiceName", "2.6.3", "983d7e1d5e8cda64")
	ctx := context.Background()
	// eventClient := public.NewStdoutClient("./AnyRobotEvent.txt")
	eventClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-983d7e1d5e8cda64/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	eventExporter := ar_event.NewExporter(eventClient)
	public.SetServiceInfo("YourServiceName", "1.0.0", "983d7e1d5e8cda64")
	eventProvider := eventsdk.NewEventProvider(eventsdk.Exporters(eventExporter, eventsdk.GetDefaultExporter()), ar_event.EventResource())
	eventsdk.SetEventProvider(eventProvider)

	defer func() {
		if err := eventProvider.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	// 业务代码
	ctx, num := add(ctx, 2, 3)
	for i := 0; i < 6; i++ {
		ctx, num = multiply(ctx, 2, 3)
		ctx, num = add(ctx, 2, 3)
	}
	log.Println(result, num)
}
