package open_standard

import (
	"span/encoder"
	"span/field"
	"time"
)

type Writer interface {
	Write(field.InternalSpan) error
}

type OpenTelemetry struct {
	Encoder encoder.Encoder
}

func (o *OpenTelemetry) Write(t field.InternalSpan) error {
	var err error
	telemetry := field.MallocStructField(10)
	telemetry.Set("TraceId", field.StringField(t.TraceID()))
	telemetry.Set("SpanId", field.StringField(t.ID()))
	telemetry.Set("ParentId", field.StringField(t.ParentID()))
	telemetry.Set("StartTime", field.TimeField(t.Time()))
	telemetry.Set("EndTime", field.TimeField(time.Now()))
	telemetry.Set("Events", field.ArrayField(t.ListRecord()))
	telemetry.Set("metrics", field.ArrayField(t.ListMetric()))
	telemetry.Set("externalSpans", field.ArrayField(t.ListExternalSpan()))
	err = o.Encoder.Write(telemetry)

	// for _, m := range t.ListMetric() {
	// 	err = o.Encoder.Write(&m)
	// }

	for _, c := range t.ListChildren() {
		err = o.Write(c)
	}

	return err
}
