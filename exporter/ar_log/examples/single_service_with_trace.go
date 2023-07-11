package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

func FileTraceInit() {
	public.SetServiceInfo("YourServiceName", "2.6.1", "983d7e1d5e8cda64")
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
	public.SetServiceInfo("YourServiceName", "2.6.1", "983d7e1d5e8cda64")
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
	public.SetServiceInfo("YourServiceName", "2.6.1", "983d7e1d5e8cda64")
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
	public.SetServiceInfo("YourServiceName", "2.6.1", "983d7e1d5e8cda64")
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
