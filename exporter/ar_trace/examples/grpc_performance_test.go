package examples

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestGRPCPerformance(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"grpc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTrace()
		})
	}
}
func setTrace() {
	wg := sync.WaitGroup{}
	ctx := context.Background()
	attrs := []attribute.KeyValue{
		attribute.String("job_id", "job-ea0ebc769228f873"),
	}
	jobResource := resource.NewWithAttributes("", attrs...)
	resource.Merge(jobResource, ar_trace.TraceResource())
	tempResource, err := resource.Merge(jobResource, ar_trace.TraceResource())
	if err == nil {
		jobResource = tempResource
	}
	traceExporter, _ := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("127.0.0.1:13034"))
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter, sdktrace.WithMaxQueueSize(5000000),
			sdktrace.WithBlocking(),
			sdktrace.WithMaxExportBatchSize(500),
			sdktrace.WithExportTimeout(time.Hour)),
		sdktrace.WithResource(jobResource))
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
	startT := time.Now()
	log.Println("开始时间：" + startT.GoString())
	for i := 1; i < 5000; i++ {
		wg.Add(1)
		go func() {
			for j := 1; j < 100000; j++ {
				time.Sleep(time.Microsecond * 50)
				multiply(ctx, 2, 3)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	tc := time.Since(startT)
	log.Printf("结束时间：" + tc.String())
}
