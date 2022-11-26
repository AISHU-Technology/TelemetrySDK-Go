package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/common"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/errors"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent"
	"encoding/json"
	"fmt"
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

	name := errors.ModuleName

	eventy := common.NewEvent(name)
	eventy.SetTime(time.Now())
	eventy.SetLevel(common.WARN)
	eventy.SetAttributes(common.NewAttribute("key", common.BoolValue(true)))
	eventy.SetLink(trace.SpanContext{})
	eventy.SetSubject("subject")
	eventy.SetData(996)
	eventy.SetEventType("type")
	eventy.GetEventMap()

	fmt.Println(eventy)

	events := make([]common.AREvent, 0)
	events = append(events, eventy)
	client := arevent.NewStdoutClient("./AnyRobotEvent.txt")
	_ = client.UploadEvents(ctx, events)

	bety, _ := json.Marshal(events)

	unmarshalEvents, err := common.UnmarshalEvents(bety)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("events", unmarshalEvents)
	for _, v := range unmarshalEvents {
		fmt.Println("event", v)
	}

	log.Println(result, num)

	//unmar := arevent.NewEvent("123")
	//b, _ := event.MarshalJSON()
	//_ = json.Unmarshal(b, &unmar)
	//println(unmar.GetEventType())
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	ctx := context.Background()
	header := make(map[string]string)
	header["self-defined"] = "some_header"
	client := arevent.NewHTTPClient(arevent.WithAnyRobotURL("https://a.b.c.d/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/events"),
		arevent.WithCompression(1), arevent.WithTimeout(10*time.Second), arevent.WithHeader(header),
		arevent.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	exporter := arevent.NewExporter(client)
	_ = exporter

	ctx, num := multiplyBefore(ctx, 7, 9)
	log.Println(result, num)
}
