package examples

import (
	"context"
	"log"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_metric"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const result = "the answer is"

// addBefore 计算两数之和。
func addBefore(ctx context.Context, x, y int64) (context.Context, int64) {
	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x + y
}

// add 增加了 Metric 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	attrs := []attribute.KeyValue{
		attribute.Key("用户信息").String("在线用户数"),
	}
	gauge, _ := ar_metric.Meter.Int64ObservableGauge("gauge：用户数峰值", instrument.WithUnit("1"), instrument.WithDescription("a simple gauge"))
	gaugeTest := func(ctx context.Context, obsrv metric.Observer) error {
		obsrv.ObserveInt64(gauge, 12, attrs...)
		return nil
	}
	attrs1 := []attribute.KeyValue{
		attribute.Key("息").String("数"),
	}
	gaugeTest1 := func(ctx context.Context, obsrv metric.Observer) error {
		obsrv.ObserveInt64(gauge, 13, attrs1...)
		return nil
	}
	_, _ = ar_metric.Meter.RegisterCallback(gaugeTest, gauge)
	_, _ = ar_metric.Meter.RegisterCallback(gaugeTest1, gauge)

	counter, _ := ar_metric.Meter.Int64ObservableCounter("CounterTest", instrument.WithUnit("1"), instrument.WithDescription("a simple gauge"))
	CounterTest := func(ctx context.Context, obsrv metric.Observer) error {
		obsrv.ObserveInt64(counter, 2, attrs...)
		return nil
	}
	_, _ = ar_metric.Meter.RegisterCallback(CounterTest, counter)

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

// multiply 增加了 Metric 的计算两数之积。
func multiply(ctx context.Context, x, y int64) (context.Context, int64) {
	attrs := []attribute.KeyValue{
		attribute.Key("用户信息").StringSlice([]string{"在线用户数"}),
	}
	histogram, _ := ar_metric.Meter.Float64Histogram("histogram：当前用户数", instrument.WithUnit("1"), instrument.WithDescription("a histogram with custom buckets and name"))
	histogram.Record(ctx, 136, attrs...)
	histogram.Record(ctx, 64, attrs...)
	histogram.Record(ctx, 340, attrs...)
	histogram.Record(ctx, 600, attrs...)

	attrs = []attribute.KeyValue{
		attribute.Key("用户信息").String("登录DAU"),
	}
	sum, _ := ar_metric.Meter.Float64Counter("sum：用户数日活", instrument.WithUnit("ms"), instrument.WithDescription("a simple counter"))
	ar_metric.Meter.Float64Counter("sum：用户数日活", instrument.WithUnit("ms"), instrument.WithDescription("a simple counter12"))
	sum.Add(ctx, 25, attrs...)
	sum.Add(ctx, 315, attrs...)
	// attrs2 := []attribute.KeyValue{
	// 	attribute.Key("用户信息").String("登录DU"),
	// }
	sum.Add(ctx, 628, attrs...)

	//业务代码
	time.Sleep(100 * time.Millisecond)
	return ctx, x * y
}

// Example 原始的业务系统入口
func Example() {
	ctx := context.Background()
	ctx, num := multiplyBefore(ctx, 2, 3)
	ctx, num = multiplyBefore(ctx, num, 7)
	ctx, num = addBefore(ctx, num, 8)
	log.Println(result, num)
}

// StdoutExample 输出到控制台和本地文件。
func StdoutExample() {
	ctx := context.Background()
	metricClient := public.NewStdoutClient("./AnyRobotMetric.txt")
	metricExporter := ar_metric.NewExporter(metricClient)
	public.SetServiceInfo("YourServiceName", "1.0.0", "")
	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second), sdkmetric.WithTimeout(10*time.Second))),
		sdkmetric.WithReader(sdkmetric.NewManualReader()),

		sdkmetric.WithResource(ar_metric.MetricResource()),
	)
	defer func() {
		if err := metricProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
	ar_metric.Meter = metricProvider.Meter(version.MetricInstrumentationName, metric.WithSchemaURL(version.MetricInstrumentationURL), metric.WithInstrumentationVersion(version.MetricInstrumentationVersion))

	ctx, num := add(ctx, 2, 8)
	ctx, num = add(ctx, 2, 8)
	ctx, num = multiply(ctx, num, 7)
	log.Println(result, num)
}

// HTTPExample 通过HTTP发送器上报到接收器。
func HTTPExample() {
	ctx := context.Background()
	metricClient := public.NewHTTPClient(public.WithAnyRobotURL("http://127.0.0.1:8800/api/feed_ingester/v1/jobs/job-abcd4f634e80d530/metrics"),
		public.WithCompression(1), public.WithTimeout(10*time.Second), public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
	metricExporter := ar_metric.NewExporter(metricClient)
	public.SetServiceInfo("YourServiceName", "1.0.0", "")
	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second), sdkmetric.WithTimeout(10*time.Second))),
		sdkmetric.WithResource(ar_metric.MetricResource()),
	)
	defer func() {
		if err := metricProvider.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
	ar_metric.Meter = metricProvider.Meter(version.MetricInstrumentationName, metric.WithSchemaURL(version.MetricInstrumentationURL), metric.WithInstrumentationVersion(version.MetricInstrumentationVersion))

	ctx, num := add(ctx, 2, 8)
	ctx, num = multiply(ctx, num, 7)
	log.Println(result, num)
}
