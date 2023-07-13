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
	"log"
	"time"
)

// SystemLogger 程序日志记录器，使用异步发送模式，无返回值。
var SystemLogger = spanLog.NewSamplerLogger(spanLog.WithSample(1.0), spanLog.WithLevel(spanLog.InfoLevel))

// ServiceLogger 业务日志记录器，使用同步发送模式，有返回值，返回error=nil代表发送成功，返回error!=nil代表发送失败。
var ServiceLogger = spanLog.NewSyncLogger(spanLog.WithLevel(spanLog.AllLevel))

const result = "the answer is"

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

// multiply 增加了 Log 的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	ctx, _ = ar_trace.Tracer.Start(ctx, "乘法")

	numbers := []int64{x, y}
	attr := field.NewAttribute("INTARRAY", field.MallocJsonField(numbers))
	// ServiceLogger 记录业务日志，并且记录自定义的业务属性信息。
	if err := ServiceLogger.Info("multiply two numbers", field.WithAttribute(attr), field.WithContext(ctx)); err != nil {
		fmt.Println(err)
	}
	// SystemLogger 记录系统日志，同时关联Log。
	SystemLogger.Fatal("This is an fatal message", field.WithContext(ctx))

	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

func FileLogInit() {
	public.SetServiceInfo("YourServiceName", "2.6.2", "983d7e1d5e8cda64")
	// 1.初始化系统日志器，系统日志写入文件。
	systemLogClient := public.NewFileClient("./AnyRobotLog.json")
	systemLogExporter := ar_log.NewExporter(systemLogClient)
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅在控制台输出。
	serviceLogExporter := exporter.SyncRealTimeExporter()
	serviceLogWriter := open_standard.NewSyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)
}

func ConsoleLogInit() {
	public.SetServiceInfo("YourServiceName", "2.6.2", "983d7e1d5e8cda64")
	// 1.初始化系统日志器，系统日志在控制台输出。
	systemLogExporter := exporter.GetRealTimeExporter()
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅在控制台输出。
	serviceLogExporter := exporter.SyncRealTimeExporter()
	serviceLogWriter := open_standard.NewSyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)
}

func HTTPLogInit() {
	public.SetServiceInfo("YourServiceName", "2.6.2", "983d7e1d5e8cda64")
	// 1.初始化系统日志器，系统日志上报AnyRobot。
	systemLogClient := public.NewHTTPClient(
		public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-864ab9d78f6a1843/events"),
		public.WithCompression(0),
		public.WithTimeout(10*time.Second),
		public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute),
	)
	systemLogExporter := ar_log.NewExporter(systemLogClient)
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅在控制台输出。
	serviceLogExporter := exporter.SyncRealTimeExporter()
	serviceLogWriter := open_standard.NewSyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)
}

func LoggerExit() {
	SystemLogger.Close()
	ServiceLogger.Close()
}

// FileExample 输出到本地文件。
func FileExample() {
	FileLogInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	LoggerExit()
	log.Println(result, num)
}

// ConsoleExample 输出到控制台。
func ConsoleExample() {
	ConsoleLogInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	LoggerExit()
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器上报到接收器。
func HTTPExample() {
	HTTPLogInit()
	ctx := context.Background()
	// 业务代码
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	LoggerExit()
	log.Println(result, num)
}

// WithAllExample 修改client所有入参。
func WithAllExample() {
	public.SetServiceInfo("YourServiceName", "2.6.2", "983d7e1d5e8cda64")
	consoleExporter := exporter.GetRealTimeExporter()

	// 1.初始化系统日志器，系统日志在控制台输出，同时上报到AnyRobot。
	systemLogClient := public.NewHTTPClient(
		public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-864ab9d78f6a1843/events"),
		public.WithCompression(0),
		public.WithTimeout(10*time.Second),
		public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute),
	)
	systemLogExporter := ar_log.NewExporter(systemLogClient)
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter, consoleExporter),
		resource.LogResource())
	systemLogRunner := runtime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	defer SystemLogger.Close()
	SystemLogger.SetLevel(spanLog.InfoLevel)
	SystemLogger.SetRuntime(systemLogRunner)

	// 2.初始化业务日志器，业务日志仅上报到AnyRobot，上报地址不同。
	serviceLogClient := public.NewSyncHTTPClient(
		public.WithAnyRobotURL("http://127.0.0.1/api/feed_ingester/v1/jobs/job-c9a577c302505576/events"),
		public.WithCompression(0),
		public.WithTimeout(10*time.Second),
		public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute),
	)
	serviceLogExporter := ar_log.NewSyncExporter(serviceLogClient)
	serviceLogWriter := open_standard.NewSyncWriter(
		encoder.NewSyncEncoder(serviceLogExporter),
		resource.LogResource())
	// 运行ServiceLogger日志器。
	defer ServiceLogger.Close()
	ServiceLogger.SetLevel(spanLog.AllLevel)
	ServiceLogger.SetWriter(serviceLogWriter)

	// 3.运行业务代码
	ctx := context.Background()
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	ctx, num := multiply(ctx, 2, 3)
	ctx, num = multiply(ctx, num, 7)
	ctx, num = add(ctx, num, 8)
	log.Println(result, num)
}
