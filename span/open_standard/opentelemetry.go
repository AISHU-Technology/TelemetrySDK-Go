package open_standard

import (
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/encoder"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
)

const rootSpan = iota

type Writer interface {
	Write([]field.LogSpan) error
	Close() error
}

type OpenTelemetry struct {
	Encoder  encoder.Encoder
	Resource field.Field
}

func NewOpenTelemetry(enc encoder.Encoder, resources field.Field) *OpenTelemetry {
	open := &OpenTelemetry{
		Encoder:  enc,
		Resource: resources,
	}
	return open
}

func (o *OpenTelemetry) Write(t []field.LogSpan) error {
	return o.write(t, rootSpan)
}

func (o *OpenTelemetry) Close() error {
	return o.Encoder.Close()
}

func (o *OpenTelemetry) write(logSpans []field.LogSpan, flag int) error {
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
