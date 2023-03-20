package open_standard

import (
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
)

const rootSpan = iota

// Writer 从自定义的runtime中获取并发送日志的写入器抽象类。
type Writer interface {
	// Write 写日志。
	Write([]field.LogSpan) error
	// Close 关闭 Writer。
	Close() error
}

// OpenTelemetry 实现了Writer抽象类的OpenTelemetry规范的日志写入器。
type OpenTelemetry struct {
	// Encoder 包含日志上报地址的编码器。
	Encoder encoder.Encoder
	// Resource 统一了数据模型后的资源信息。
	Resource field.Field
}

// NewOpenTelemetry 对外暴露的由客户调用的初始化OpenTelemetry规范的日志写入器的方法。
func NewOpenTelemetry(enc encoder.Encoder, resources field.Field) *OpenTelemetry {
	open := &OpenTelemetry{
		Encoder:  enc,
		Resource: resources,
	}
	return open
}

func (o *OpenTelemetry) Write(t []field.LogSpan) error {
	return o.write(t)
}

func (o *OpenTelemetry) Close() error {
	return o.Encoder.Close()
}

func (o *OpenTelemetry) write(logSpans []field.LogSpan) error {
	var err error
	telemetrys := field.MallocArrayField(len(logSpans) + 1)
	for _, t := range logSpans {
		telemetry := field.MallocStructField(8)

		link := field.MallocStructField(2)
		link.Set("TraceId", field.StringField(t.TraceID()))
		link.Set("SpanId", field.StringField(t.SpanID()))

		telemetry.Set("Link", link)
		telemetry.Set("Timestamp", field.StringField(time.Now().Format(time.RFC3339Nano)))
		telemetry.Set("SeverityText", t.GetLogLevel())

		telemetry.Set("Body", t.GetRecord())
		attrs := t.GetAttributes()

		telemetry.Set("Attributes", attrs)
		telemetry.Set("Resource", o.Resource)
		telemetrys.Append(telemetry)
	}

	err = o.Encoder.Write(telemetrys)
	if err != nil {
		return err
	}

	return err
}
