package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_metric"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
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

// add 增加了 Metric 的计算两数之和。
func add(ctx context.Context, x, y int64) (context.Context, int64) {
	attrs := []attribute.KeyValue{
		attribute.Key("keyA").String("valueB"),
		attribute.Key("keyC").StringSlice([]string{"valueD", "valueE"}),
	}
	gauge, _ := ar_metric.Meter.AsyncInt64().Gauge("test_gauge")
	gaugeTest := func(ctx context.Context) {
		gauge.Observe(ctx, 3.0, attrs...)
	}
	_ = ar_metric.Meter.RegisterCallback([]instrument.Asynchronous{gauge}, gaugeTest)

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
		attribute.Key("keyA").String("valueB"),
		attribute.Key("keyC").StringSlice([]string{"valueD", "valueE"}),
	}
	histogram, _ := ar_metric.Meter.SyncFloat64().Histogram("test_histogram", instrument.WithUnit(unit.Dimensionless), instrument.WithDescription("a histogram with custom buckets and rename"))
	histogram.Record(ctx, 136, attrs...)
	histogram.Record(ctx, 64, attrs...)

	sum, _ := ar_metric.Meter.SyncFloat64().Counter("test_sum", instrument.WithUnit(unit.Milliseconds), instrument.WithDescription("a simple counter"))
	sum.Add(ctx, 5, attrs...)
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

// HTTPExample 通过HTTP发送器上报到接收器。
func HTTPExample() {
	ctx := context.Background()

	ctx, num := multiply(ctx, 2, 3)
	for i := 0; i < 6; i++ {
		ctx, num = multiply(ctx, 2, 3)
		ctx, num = add(ctx, 2, 3)
	}
	log.Println(result, num)
}
