package open_standard

import (
    "os"
    "span/encoder"
    "span/field"
    "time"
)

const (
    rootSpan = iota
)

type Writer interface {
    Write(field.InternalSpan) error
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

func (o *OpenTelemetry) Write(t field.InternalSpan) error {
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
    f.Set("telemetry.sdk.name", field.StringField("Aishu custom opentelemetry"))
    f.Set("telemetry.sdk.version", field.StringField("1.0.0"))
    f.Set("telemetry.sdk.language", field.StringField("go"))

    o.Resource = f
}

func (o *OpenTelemetry) Close() error {
    return o.Encoder.Close()
}

func (o *OpenTelemetry) write(t field.InternalSpan, flag int) error {
    var err error
    telemetry := field.MallocStructField(8)
    telemetry.Set("Version", field.StringField("AISHUV0"))
    telemetry.Set("TraceId", field.StringField(t.TraceID()))
    telemetry.Set("SpanId", field.StringField(t.ID()))
    telemetry.Set("ParentId", field.StringField(t.ParentID()))
    telemetry.Set("StartTime", field.TimeField(t.Time()))
    telemetry.Set("EndTime", field.TimeField(time.Now()))

    events := field.ArrayField(t.ListRecord())
    body := field.MallocStructField(3)
    body.Set("Events", &events)
    metrics := field.ArrayField(t.ListMetric())
    body.Set("Metrics", &metrics)
    external := field.ArrayField(t.ListExternalSpan())
    body.Set("ExternalSpans", &external)
    telemetry.Set("Body", body)
    attrs := t.GetAttributes()
    if attrs != nil {
        telemetry.Set("Attributes", attrs)
    }

    if o.Resource == nil {
        o.SetDefultResources()
    }

    telemetry.Set("Resource", o.Resource)

    err = o.Encoder.Write(telemetry)
    if err != nil {
        return err
    }

    // for _, m := range t.ListMetric() {
    //  err = o.Encoder.Write(&m)
    // }

    for _, c := range t.ListChildren() {
        err = o.Write(c)
    }

    return err
}

