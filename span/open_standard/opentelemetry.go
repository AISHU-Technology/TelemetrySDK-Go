package open_standard

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"

	"os"
	"time"
)

const (
	rootSpan                = iota
	OpenTelemetrySDKVersion = "v1.6.1"
	SDKName                 = "Telemetry SDK"
	SDKVersion              = "2.0.0"
	SDKLanguage             = "go"
)

type Writer interface {
	Write(field.LogSpan) error
	Close() error
}

type OpenTelemetry struct {
	Encoder  encoder.Encoder
	Resource field.Field
}

func NewOpenTelemetry(enc encoder.Encoder, resources field.Field) OpenTelemetry {
	res := OpenTelemetry{
		Encoder:  enc,
		Resource: resources,
	}
	if res.Resource == nil {
		res.SetDefultResources()
	}

	return res
}

func (o *OpenTelemetry) Write(t field.LogSpan) error {
	return o.write(t, rootSpan)
}

func (o *OpenTelemetry) SetDefultResources() {
	f := field.MallocStructField(10)
	targets := []string{"HOSTNAME"}
	for _, k := range targets {
		if v, e := os.LookupEnv(k); e {
			f.Set(k, field.StringField(v))
		}
	}
	f.Set("telemetry.sdk.name", field.StringField(SDKName))
	f.Set("telemetry.sdk.version", field.StringField(SDKVersion))
	f.Set("telemetry.sdk.language", field.StringField(SDKLanguage))

	o.Resource = f
}

func (o *OpenTelemetry) Close() error {
	return o.Encoder.Close()
}

func (o *OpenTelemetry) write(t field.LogSpan, flag int) error {
	var err error
	telemetry := field.MallocStructField(8)
	telemetry.Set("Version", field.StringField(OpenTelemetrySDKVersion))
	telemetry.Set("TraceId", field.StringField(t.TraceID()))
	telemetry.Set("SpanId", field.StringField(t.SpanID()))
	telemetry.Set("Timestamp", field.TimeField(time.Now()))

	telemetry.Set("Body", t.GetRecord())
	attrs := t.GetAttributes()

	telemetry.Set("Attributes", attrs)

	if o.Resource == nil {
		o.SetDefultResources()
	}

	telemetry.Set("Resource", o.Resource)

	err = o.Encoder.Write(telemetry)
	if err != nil {
		return err
	}

	return err
}
