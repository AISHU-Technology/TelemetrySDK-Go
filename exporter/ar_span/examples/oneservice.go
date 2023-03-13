package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_span"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	spanLog "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/log"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/open_standard"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/runtime"
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

// add 增加了 Log 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, _ = ar_trace.Tracer.Start(ctx, "加法")

	numbers := []int64{x, y}
	attr := field.NewAttribute("INTARRAY", field.MallocJsonField(numbers))
	ar_span.ServiceLogger.Info("add two numbers", field.WithAttribute(attr))
	ar_span.SystemLogger.Error("This is an error message")

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

	numbers := []int64{x, y}
	attr := field.NewAttribute("INTARRAY", field.MallocJsonField(numbers))
	ar_span.ServiceLogger.Info("multiply two numbers", field.WithAttribute(attr))
	ar_span.SystemLogger.Fatal("This is an fatal message")

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

// HTTPExample 修改client所有入参。
func HTTPExample() {
	public.SetServiceInfo("YourServiceName", "1.0.0", "")
	stdoutExporter := exporter.GetStdoutExporter()

	systemLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://10.4.130.68:13048/api/feed_ingester/v1/jobs/Kitty1/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	systemLogExporter := ar_span.NewExporter(systemLogClient)
	systemLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(systemLogExporter, stdoutExporter),
		Resource: resource.LogResource(),
	}
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)

	serviceLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://10.4.130.68:13048/api/feed_ingester/v1/jobs/Kitty1/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	serviceLogExporter := ar_span.NewExporter(serviceLogClient)
	serviceLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(serviceLogExporter, stdoutExporter),
		Resource: resource.LogResource(),
	}
	serviceLogRunner := runtime.NewRuntime(serviceLogWriter, field.NewSpanFromPool)
	serviceLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)

	// start runtime
	go systemLogRunner.Run()
	ar_span.SystemLogger.SetLevel(spanLog.AllLevel)
	ar_span.SystemLogger.SetRuntime(systemLogRunner)

	ar_span.ServiceLogger.SetLevel(spanLog.AllLevel)
	ar_span.ServiceLogger.SetRuntime(systemLogRunner)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	ar_span.SystemLogger.Close()
	ar_span.ServiceLogger.Close()
}

func StdoutExporterExample() {
	stdoutExporter := exporter.GetStdoutExporter()
	public.SetServiceInfo("YourServiceName", "1.0.0", "")
	writer := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(stdoutExporter),
		Resource: resource.LogResource(),
	}
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	run.SetUploadInternalAndMaxLog(3*time.Second, 10)

	// start runtime
	go run.Run()
	ar_span.SystemLogger.SetLevel(spanLog.AllLevel)
	ar_span.SystemLogger.SetRuntime(run)

	ar_span.ServiceLogger.SetLevel(spanLog.AllLevel)
	ar_span.ServiceLogger.SetRuntime(run)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	ar_span.SystemLogger.Close()
	ar_span.ServiceLogger.Close()
}

func OldStdoutExample() {
	public.SetServiceInfo("YourServiceName", "1.0.0", "")
	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
		//Resource: field.WithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID"),
	}
	//writer.SetDefaultResources()
	//writer.SetResourcesWithServiceInfo("testServiceName", "testServiceVersion", "testServiceInstanceID")
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)

	// start runtime
	go run.Run()
	ar_span.SystemLogger.SetLevel(spanLog.AllLevel)
	ar_span.SystemLogger.SetRuntime(run)

	ar_span.ServiceLogger.SetLevel(spanLog.AllLevel)
	ar_span.ServiceLogger.SetRuntime(run)

	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
	ar_span.SystemLogger.Close()
	ar_span.ServiceLogger.Close()
}
