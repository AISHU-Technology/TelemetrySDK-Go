package open_standard

import (
	"span/encoder"
	"span/field"
	"time"
)

const (
	rootSpan = iota
)

type Writer interface {
	Write(field.InternalSpan) error
}

type OpenTelemetry struct {
	Encoder encoder.Encoder
}

func (o *OpenTelemetry) Write(t field.InternalSpan) error {
	return o.write(t, rootSpan)
}

func (o *OpenTelemetry) write(t field.InternalSpan, flag int) error {
	var err error
	telemetry := field.MallocStructField(8)
	telemetry.Set("TraceId", field.StringField(t.TraceID()))
	telemetry.Set("SpanId", field.StringField(t.ID()))
	telemetry.Set("ParentId", field.StringField(t.ParentID()))
	telemetry.Set("StartTime", field.TimeField(t.Time()))
	telemetry.Set("EndTime", field.TimeField(time.Now()))

	events := field.ArrayField(t.ListRecord())
	telemetry.Set("Events", &events)
	metrics := field.ArrayField(t.ListMetric())
	telemetry.Set("metrics", &metrics)
	external := field.ArrayField(t.ListExternalSpan())
	telemetry.Set("externalSpans", &external)

	err = o.Encoder.Write(telemetry)
	if err != nil {
		return err
	}

	// for _, m := range t.ListMetric() {
	// 	err = o.Encoder.Write(&m)
	// }

	for _, c := range t.ListChildren() {
		err = o.Write(c)
	}

	return err
}
