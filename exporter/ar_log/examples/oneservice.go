package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_log"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/exporter"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	spanLog "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/log"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/open_standard"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/runtime"
	"fmt"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

// SystemLogger 程序日志记录器，使用异步发送模式，无返回值。
var SystemLogger = spanLog.NewSamplerLogger(spanLog.WithSample(1.0), spanLog.WithLevel(spanLog.InfoLevel))

// ServiceLogger 业务日志记录器，使用同步发送模式，有返回值，返回error=nil代表发送成功，返回error!=nil代表发送失败。
var ServiceLogger = spanLog.NewSyncLogger(spanLog.WithLevel(spanLog.AllLevel))

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
	// ServiceLogger 记录业务日志，并且记录自定义的业务属性信息，同步发送处理返回结果。
	if err := ServiceLogger.Info("add two numbers", field.WithAttribute(attr)); err != nil {
		fmt.Println(err)
	}
	// SystemLogger 记录系统日志。
	SystemLogger.Error("This is an error message")

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
	// ServiceLogger 记录业务日志，并且记录自定义的业务属性信息。
	if err := ServiceLogger.Info("multiply two numbers", field.WithAttribute(attr), field.WithContext(ctx)); err != nil {
		fmt.Println(err)
	}
	// SystemLogger 记录系统日志，同时关联Trace。
	SystemLogger.Fatal("This is an fatal message", field.WithContext(ctx))

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
	fmt.Println(result, num)
}

// HTTPExample 修改client所有入参。
func HTTPExample() {
	public.SetServiceInfo("YourServiceName", "2.6.5", "c9a577c302505576")
	stdoutExporter := exporter.GetStdoutExporter()

	// 1.初始化系统日志器，系统日志在控制台输出，同时上报到AnyRobot。
	systemLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-983d7e1d5e8cda64/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	systemLogExporter := ar_log.NewExporter(systemLogClient)
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter, stdoutExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	//defer SystemLogger.Close()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅上报到AnyRobot，上报地址不同。
	serviceLogClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-c9a577c302505576/events"),
		public.WithCompression(0), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	serviceLogExporter := ar_log.NewExporter(serviceLogClient)
	serviceLogWriter := open_standard.SyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	//defer ServiceLogger.Close()
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)

	// 3.运行业务代码
	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	_, _ = add(ctx, num, 8)
}

func StdoutExporterExample() {
	public.SetServiceInfo("YourServiceName", "2.6.5", "983d7e1d5e8cda64")

	// 1.初始化系统日志器，系统日志在控制台输出。
	systemLogExporter := exporter.GetStdoutExporter()
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	//defer SystemLogger.Close()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅在控制台输出。
	serviceLogExporter := exporter.GetStdoutExporter()
	serviceLogWriter := open_standard.SyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	//defer ServiceLogger.Close()
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)

	// 3.运行业务代码
	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 1, 7)
	ctx, _ = add(ctx, num, 8)
	for i := 0; i < 10; i++ {
		ctx, num = multiply(ctx, 1, num)
	}
}
