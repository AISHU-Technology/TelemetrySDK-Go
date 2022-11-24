package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
	"time"
)

const result = "the answer is"

// addBefore 计算两数之和。
func addBefore(ctx context.Context, x, y int64) (context.Context, int64) {
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

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	ctx := context.Background()
	//client := arevent.NewStdoutClient("./AnyRobotEvent.txt")

	ctx, num := addBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)

	event := arevent.NewEvent("examples.exporters.arevent")
	event.SetSubject("stdout.example")
	event.SetLevel(arevent.WARN)
	event.SetAttributes(arevent.NewAttribute("样例", arevent.StringValue("结果")))
	event.SetData(num)
	tracerProvider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	_, span := artrace.Tracer.Start(ctx, "")
	event.SetLink(span.SpanContext())

	event.SetAttributes(arevent.NewAttribute("", arevent.StringValue("123")))
	//_ = client.UploadEvent(ctx, &event)

	event3 := arevent.NewEvent("")
	mapping := event3.GetEventMap()
	fmt.Println(mapping)
	file1 := os.Stdout
	encoder1 := json.NewEncoder(file1)
	encoder1.SetEscapeHTML(false)
	encoder1.SetIndent("", "\t")
	_ = encoder1.Encode(mapping)

	log.Println(result, num)
}
