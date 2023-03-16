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
	// ar_span.ServiceLogger 记录业务日志，并且记录自定义的业务属性信息。
	ar_span.ServiceLogger.Info("add two numbers", field.WithAttribute(attr))
	// ar_span.SystemLogger 记录系统日志。
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

// multiply 增加了 Log 的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, _ = ar_trace.Tracer.Start(ctx, "乘法")

	numbers := []int64{x, y}
	attr := field.NewAttribute("INTARRAY", field.MallocJsonField(numbers))
	// ar_span.ServiceLogger 记录业务日志，并且记录自定义的业务属性信息。
	ar_span.ServiceLogger.Info("multiply two numbers", field.WithAttribute(attr), field.WithContext(ctx))
	// ar_span.SystemLogger 记录系统日志，同时关联Trace。
	ar_span.SystemLogger.Fatal("This is an fatal message", field.WithContext(ctx))

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

	// 1.初始化系统日志器，系统日志在控制台输出，同时上报到AnyRobot。
	systemLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1:8800/api/feed_ingester/v1/jobs/Kitty1/events1"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	systemLogExporter := ar_span.NewExporter(systemLogClient)
	systemLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(systemLogExporter, stdoutExporter),
		Resource: resource.LogResource(),
	}
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	defer ar_span.SystemLogger.Close()
	ar_span.SystemLogger.SetLevel(spanLog.InfoLevel)
	ar_span.SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅上报到AnyRobot，上报地址不同。
	serviceLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1:9900/api/feed_ingester/v1/jobs/Kitty2/events2"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	serviceLogExporter := ar_span.NewExporter(serviceLogClient)
	serviceLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(serviceLogExporter),
		Resource: resource.LogResource(),
	}
	serviceLogRunner := runtime.NewRuntime(serviceLogWriter, field.NewSpanFromPool)
	serviceLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行ServiceLogger日志器。
	go serviceLogRunner.Run()
	defer ar_span.ServiceLogger.Close()
	ar_span.ServiceLogger.SetLevel(spanLog.AllLevel)
	ar_span.ServiceLogger.SetRuntime(serviceLogRunner)

	// 3.运行业务代码
	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)

}

func StdoutExporterExample() {
	public.SetServiceInfo("YourServiceName", "1.0.0", "")

	// 1.初始化系统日志器，系统日志在控制台输出。
	systemLogExporter := exporter.GetStdoutExporter()
	systemLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(systemLogExporter),
		Resource: resource.LogResource(),
	}
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	defer ar_span.SystemLogger.Close()
	ar_span.SystemLogger.SetLevel(spanLog.InfoLevel)
	ar_span.SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅在控制台输出。
	serviceLogExporter := exporter.GetStdoutExporter()
	serviceLogWriter := &open_standard.OpenTelemetry{
		Encoder:  encoder.NewJsonEncoderWithExporters(serviceLogExporter),
		Resource: resource.LogResource(),
	}
	serviceLogRunner := runtime.NewRuntime(serviceLogWriter, field.NewSpanFromPool)
	serviceLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行ServiceLogger日志器。
	go serviceLogRunner.Run()
	defer ar_span.ServiceLogger.Close()
	ar_span.ServiceLogger.SetLevel(spanLog.AllLevel)
	ar_span.ServiceLogger.SetRuntime(serviceLogRunner)

	// 3.运行业务代码
	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
}
