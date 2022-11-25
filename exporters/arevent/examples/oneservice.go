package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"
	"encoding/json"
	"fmt"
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

// multiplyBefore 计算两数之积。
func multiplyBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	ctx := context.Background()

	ctx, num := addBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)

	event := arevent.NewEvent("examples.exporters.arevent")
	event2 := arevent.NewEvent("examples.exporters.arevent2")
	event3 := arevent.NewEvent("examples.exporters.arevent3")
	//println(event.GetEventMap())

	//file1 := os.Stdout
	//encoder1 := json.NewEncoder(file1)
	//encoder1.SetEscapeHTML(false)
	//encoder1.SetIndent("", "\t")
	//_ = encoder1.Encode(event)

	events := make([]model.AREvent, 0)
	events = append(events, event)
	events = append(events, event2)
	events = append(events, event3)
	//client := arevent.NewStdoutClient("./AnyRobotEvent.txt")
	//_ = client.UploadEvents(ctx, events)

	bety, _ := json.Marshal(events)

	unmarshalEvents, err := arevent.UnmarshalEvents(bety)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range unmarshalEvents {
		fmt.Println(v.GetEventID())
	}

	//betyy, _ := json.Marshal(event.GetEventType())

	//fmt.Println(unmarshalEvents)
	//

	//file1 := os.Stdout
	//encoder1 := json.NewEncoder(file1)
	//encoder1.SetEscapeHTML(false)
	//encoder1.SetIndent("", "\t")
	//_ = encoder1.Encode(unmarshalEvents)

	//event.SetSubject("stdout.example")
	//event.SetLevel(arevent.WARN)
	//event.SetAttributes(arevent.NewAttribute("样例", arevent.StringValue("结果")))
	//event.SetData(num)
	//tracerProvider := sdktrace.NewTracerProvider()
	//otel.SetTracerProvider(tracerProvider)
	//_, span := artrace.Tracer.Start(ctx, "")
	//event.SetLink(span.SpanContext())

	//event.SetAttributes(arevent.NewAttribute("", arevent.StringValue("123")))
	//_ = client.UploadEvent(context.Background(), &event)

	//event3 := arevent.NewEvent("123")
	//mapping := event3.GetEventMap()
	//fmt.Println(mapping)
	//file1 := os.Stdout
	//encoder1 := json.NewEncoder(file1)
	//encoder1.SetEscapeHTML(false)
	//encoder1.SetIndent("", "\t")
	//_ = encoder1.Encode(mapping)

	//events := make([]model.AREvent, 0)
	//events = append(events, event)
	//events = append(events, event3)
	//bety, _ := json.Marshal(events)
	//unmarshalEvents, err := arevent.UnmarshalEvents(bety)
	//results := make([]model.AREvent, 0)
	//err := json.Unmarshal(bety, &results)
	//if err != nil {
	//	println(err)
	//}

	//fmt.Println(events)

	//println(results)
	//client := arevent.NewStdoutClient("")
	//client.UploadEvent(context.Background(), &event3)
	//file1 := os.Stdout
	//encoder1 := json.NewEncoder(file1)
	//encoder1.SetEscapeHTML(false)
	//encoder1.SetIndent("", "\t")
	//_ = encoder1.Encode(&event3)

	log.Println(result, num)

	//unmar := arevent.NewEvent("123")
	//b, _ := event.MarshalJSON()
	//_ = json.Unmarshal(b, &unmar)
	//println(unmar.GetEventType())
}
