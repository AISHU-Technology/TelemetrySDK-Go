package examples

import (
	"context"
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
	time.Sleep(300 * time.Millisecond)
	return ctx, x + y
}

// multiplyBefore 计算两数之积。
func multiplyBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(300 * time.Millisecond)
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
	client := arevent.NewStdoutClient("./AnyRobotEvent.txt")

	ctx, num := multiplyBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)
	event := arevent.NewEvent("34")
	event.SetEventType("com.github.pull.create")
	event.SetTime(time.Now())
	event.SetLevel(arevent.ERROR)
	tracerProvider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	_, span := artrace.Tracer.Start(ctx, "加法", trace.WithSpanKind(1))
	event.SetLink(span.SpanContext())
	event.SetAttributes(arevent.NewAttribute("123", arevent.BoolValue(false)))

	event.SetAttributes(arevent.NewAttribute("key", arevent.FloatArray([]float64{1.2, 3.4, 5.6})))
	event.SetSubject("sth.png")
	event.SetData(987456)

	_ = client.UploadEvent(ctx, &event)

	event2 := arevent.NewEvent("")
	event2.SetTime(time.UnixMicro(1212654612121))
	event2.SetData(&event)
	_ = client.UploadEvent(ctx, &event2)
	ctx, num = addBefore(ctx, num, 8)
	log.Println(result, num)

	//time.RFC3339Nano
}
