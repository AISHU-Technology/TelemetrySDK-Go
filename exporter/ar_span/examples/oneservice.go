package examples

import (
	"context"
	"log"
	"os"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_span"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	spanLog "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/log"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/open_standard"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/runtime"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var sampleLog *spanLog.SamplerLogger

type LogModel struct {
	Model   string
	Project string
}

var lm = LogModel{
	Model:   "main",
	Project: "test-demo",
}

const result = "the answer is"

// addBefore 计算两数之和。
func addBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// add 增加了 Log 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, _ = ar_trace.Tracer.Start(ctx, "加法")
	attr := field.NewAttribute("test", field.StringField("test"))
	//结构化日志
	sampleLog.InfoField(field.MallocJsonField(lm), "test", field.WithContext(ctx), field.WithAttribute(attr))
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
	ctx, _ = ar_trace.Tracer.Start(ctx, "乘法")
	attr := field.NewAttribute("logModel", field.MallocJsonField(lm))
	// 设置 attribute的日志
	sampleLog.Info("this is a test multiply", field.WithContext(ctx), field.WithAttribute(attr))
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y

}

// Example 原始的业务系统入口
func Example() {
	ctx := context.Background()
	ctx, num := multiplyBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)
	_, num = addBefore(ctx, num, 8)
	log.Println(result, num)
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExporterExample() {
	sampleLog = spanLog.NewDefaultSamplerLogger()
	// // 此处demo output为标准输出
	stdoutLogClient := public.NewStdoutClient("./AnyRobotLog.txt")
	stdoutLogExporter := ar_span.NewExporter(stdoutLogClient)
	writer := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(stdoutLogExporter),
		Resource: field.WithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID"),
	}
	//writer.SetDefaultResources()
	//writer.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)

	run.SetUploadInternalAndMaxLog(3*time.Second, 10)

	sampleLog.SetRuntime(run)

	// start runtime
	go run.Run()
	sampleLog.SetLevel(spanLog.AllLevel)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	sampleLog.Info("this is a test", field.WithContext(ctx))
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	sampleLog.Close()
}

// HTTPExample 修改client所有入参。
func HTTPExample() {
	sampleLog = spanLog.NewDefaultSamplerLogger()
	defaultExporter := exporter.GetDefaultExporter()
	logClient := public.NewHTTPClient(public.WithAnyRobotURL("http://10.4.130.68:13048/api/feed_ingester/v1/jobs/Kitty1/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	logExporter := ar_span.NewExporter(logClient)
	//defaultExporter := exporter.GetDefaultExporter()
	//	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoderWithExporters(logExporter, defaultExporter),
		//Resource: field.WithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID"),
	}
	//writer.SetDefaultResources()
	writer.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)

	run.SetUploadInternalAndMaxLog(3*time.Second, 10)

	sampleLog.SetRuntime(run)

	// start runtime
	go run.Run()
	sampleLog.SetLevel(spanLog.AllLevel)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	sampleLog.Info("this is a test", field.WithContext(ctx))
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	sampleLog.Close()
}

func DefaultExporterExample() {
	sampleLog = spanLog.NewDefaultSamplerLogger()
	defaultExporter := exporter.GetDefaultExporter()
	//	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoderWithExporters(defaultExporter),
		//Resource: field.WithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID"),
	}
	writer.SetDefaultResources()
	//writer.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)

	run.SetUploadInternalAndMaxLog(3*time.Second, 10)

	sampleLog.SetRuntime(run)

	// start runtime
	go run.Run()
	sampleLog.SetLevel(spanLog.AllLevel)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	sampleLog.Info("this is a test", field.WithContext(ctx))
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	sampleLog.Close()
}

func StdoutExample() {
	sampleLog = spanLog.NewDefaultSamplerLogger()
	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
		//Resource: field.WithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID"),
	}
	//writer.SetDefaultResources()
	//writer.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)

	sampleLog.SetRuntime(run)

	// start runtime
	go run.Run()
	sampleLog.SetLevel(spanLog.AllLevel)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	sampleLog.Info("this is a test", field.WithContext(ctx))
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	sampleLog.Close()
}
