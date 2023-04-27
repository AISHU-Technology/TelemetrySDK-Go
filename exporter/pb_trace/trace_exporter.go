package pb_trace

import (
	"context"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/version"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

// 跨包实现接口占位用。
var _ sdktrace.SpanExporter = (*TraceExporter)(nil)

// Tracer 是一个全局变量，用于在业务代码中生产Span。
var Tracer = otel.GetTracerProvider().Tracer(
	version.TraceInstrumentationName,
	trace.WithInstrumentationVersion(version.TelemetrySDKVersion),
	trace.WithSchemaURL(version.TraceInstrumentationURL),
)

// TraceExporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type TraceExporter struct {
	*public.Exporter
}

// ExportSpans 批量发送AnyRobotSpans到AnyRobot Feed Ingester的Trace数据接收器。
func (e *TraceExporter) ExportSpans(ctx context.Context, traces []sdktrace.ReadOnlySpan) error {
	if len(traces) == 0 {
		return nil
	}
	protoSpans := Spans(traces)
	for _, protoSpan := range protoSpans {
		sendBytes, err := proto.Marshal(protoSpan)
		if err != nil {
			err=e.ExportData(ctx, sendBytes)
			if err != nil {
				return err
			}
		}else{
			return err
		}
	}
	return nil
	//return e.ExportData(ctx, file.Bytes())
}

// NewExporter 创建已启动的 TraceExporter 。
func NewExporter(c public.Client) *TraceExporter {
	return &TraceExporter{
		public.NewExporter(c),
	}
}

// TraceResource 传入 Trace 的默认Resource。
func TraceResource() *sdkresource.Resource {
	return resource.TraceResource()
}

func Spans(sdl []sdktrace.ReadOnlySpan) []*tracepb.ResourceSpans {
	if len(sdl) == 0 {
		return nil
	}

	rsm := make(map[attribute.Distinct]*tracepb.ResourceSpans)

	type key struct {
		r  attribute.Distinct
		is instrumentation.Scope
	}
	ssm := make(map[key]*tracepb.ScopeSpans)

	var resources int
	for _, sd := range sdl {
		if sd == nil {
			continue
		}

		rKey := sd.Resource().Equivalent()
		k := key{
			r:  rKey,
			is: sd.InstrumentationScope(),
		}
		scopeSpan, iOk := ssm[k]
		if !iOk {
			// Either the resource or instrumentation scope were unknown.
			scopeSpan = &tracepb.ScopeSpans{
				Scope:     InstrumentationScope(sd.InstrumentationScope()),
				Spans:     []*tracepb.Span{},
				SchemaUrl: sd.InstrumentationScope().SchemaURL,
			}
		}
		scopeSpan.Spans = append(scopeSpan.Spans, span(sd))
		ssm[k] = scopeSpan

		rs, rOk := rsm[rKey]
		if !rOk {
			resources++
			// The resource was unknown.
			rs = &tracepb.ResourceSpans{
				Resource:   Resource(sd.Resource()),
				ScopeSpans: []*tracepb.ScopeSpans{scopeSpan},
				SchemaUrl:  sd.Resource().SchemaURL(),
			}
			rsm[rKey] = rs
			continue
		}

		// The resource has been seen before. Check if the instrumentation
		// library lookup was unknown because if so we need to add it to the
		// ResourceSpans. Otherwise, the instrumentation library has already
		// been seen and the append we did above will be included it in the
		// ScopeSpans reference.
		if !iOk {
			rs.ScopeSpans = append(rs.ScopeSpans, scopeSpan)
		}
	}

	// Transform the categorized map into a slice
	rss := make([]*tracepb.ResourceSpans, 0, resources)
	for _, rs := range rsm {
		rss = append(rss, rs)
	}
	return rss
}

// Resource transforms a Resource into an OTLP Resource.
func Resource(r *sdkresource.Resource) *resourcepb.Resource {
	if r == nil {
		return nil
	}
	return &resourcepb.Resource{Attributes: ResourceAttributes(r)}
}

// ResourceAttributes transforms a Resource OTLP key-values.
func ResourceAttributes(res *sdkresource.Resource) []*commonpb.KeyValue {
	return Iterator(res.Iter())
}

// Iterator transforms an attribute iterator into OTLP key-values.
func Iterator(iter attribute.Iterator) []*commonpb.KeyValue {
	l := iter.Len()
	if l == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, l)
	for iter.Next() {
		out = append(out, KeyValue(iter.Attribute()))
	}
	return out
}

// span transforms a Span into an OTLP span.
func span(sd sdktrace.ReadOnlySpan) *tracepb.Span {
	if sd == nil {
		return nil
	}

	tid := sd.SpanContext().TraceID()
	sid := sd.SpanContext().SpanID()

	s := &tracepb.Span{
		TraceId:                tid[:],
		SpanId:                 sid[:],
		TraceState:             sd.SpanContext().TraceState().String(),
		Status:                 status(sd.Status().Code, sd.Status().Description),
		StartTimeUnixNano:      uint64(sd.StartTime().UnixNano()),
		EndTimeUnixNano:        uint64(sd.EndTime().UnixNano()),
		Links:                  links(sd.Links()),
		Kind:                   spanKind(sd.SpanKind()),
		Name:                   sd.Name(),
		Attributes:             KeyValues(sd.Attributes()),
		Events:                 spanEvents(sd.Events()),
		DroppedAttributesCount: uint32(sd.DroppedAttributes()),
		DroppedEventsCount:     uint32(sd.DroppedEvents()),
		DroppedLinksCount:      uint32(sd.DroppedLinks()),
	}

	if psid := sd.Parent().SpanID(); psid.IsValid() {
		s.ParentSpanId = psid[:]
	}

	return s
}

// spanEvents transforms span Events to an OTLP span events.
func spanEvents(es []sdktrace.Event) []*tracepb.Span_Event {
	if len(es) == 0 {
		return nil
	}

	events := make([]*tracepb.Span_Event, len(es))
	// Transform message events
	for i := 0; i < len(es); i++ {
		events[i] = &tracepb.Span_Event{
			Name:                   es[i].Name,
			TimeUnixNano:           uint64(es[i].Time.UnixNano()),
			Attributes:             KeyValues(es[i].Attributes),
			DroppedAttributesCount: uint32(es[i].DroppedAttributeCount),
		}
	}
	return events
}

// spanKind transforms a SpanKind to an OTLP span kind.
func spanKind(kind trace.SpanKind) tracepb.Span_SpanKind {
	switch kind {
	case trace.SpanKindInternal:
		return tracepb.Span_SPAN_KIND_INTERNAL
	case trace.SpanKindClient:
		return tracepb.Span_SPAN_KIND_CLIENT
	case trace.SpanKindServer:
		return tracepb.Span_SPAN_KIND_SERVER
	case trace.SpanKindProducer:
		return tracepb.Span_SPAN_KIND_PRODUCER
	case trace.SpanKindConsumer:
		return tracepb.Span_SPAN_KIND_CONSUMER
	default:
		return tracepb.Span_SPAN_KIND_UNSPECIFIED
	}
}

// status transform a span code and message into an OTLP span status.
func status(status codes.Code, message string) *tracepb.Status {
	var c tracepb.Status_StatusCode
	switch status {
	case codes.Ok:
		c = tracepb.Status_STATUS_CODE_OK
	case codes.Error:
		c = tracepb.Status_STATUS_CODE_ERROR
	default:
		c = tracepb.Status_STATUS_CODE_UNSET
	}
	return &tracepb.Status{
		Code:    c,
		Message: message,
	}
}

func links(links []sdktrace.Link) []*tracepb.Span_Link {
	if len(links) == 0 {
		return nil
	}

	sl := make([]*tracepb.Span_Link, 0, len(links))
	for _, otLink := range links {
		// This redefinition is necessary to prevent otLink.*ID[:] copies
		// being reused -- in short we need a new otLink per iteration.
		otLink := otLink

		tid := otLink.SpanContext.TraceID()
		sid := otLink.SpanContext.SpanID()

		sl = append(sl, &tracepb.Span_Link{
			TraceId:                tid[:],
			SpanId:                 sid[:],
			Attributes:             KeyValues(otLink.Attributes),
			DroppedAttributesCount: uint32(otLink.DroppedAttributeCount),
		})
	}
	return sl
}

// KeyValues transforms a slice of attribute KeyValues into OTLP key-values.
func KeyValues(attrs []attribute.KeyValue) []*commonpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, KeyValue(kv))
	}
	return out
}

// KeyValue transforms an attribute KeyValue into an OTLP key-value.
func KeyValue(kv attribute.KeyValue) *commonpb.KeyValue {
	return &commonpb.KeyValue{Key: string(kv.Key), Value: Value(kv.Value)}
}

// Value transforms an attribute Value into an OTLP AnyValue.
func Value(v attribute.Value) *commonpb.AnyValue {
	av := new(commonpb.AnyValue)
	switch v.Type() {
	case attribute.BOOL:
		av.Value = &commonpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case attribute.BOOLSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		av.Value = &commonpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.INT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		av.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.FLOAT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.STRINGSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	default:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func boolSliceValues(vals []bool) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_BoolValue{
				BoolValue: v,
			},
		}
	}
	return converted
}

func int64SliceValues(vals []int64) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: v,
			},
		}
	}
	return converted
}

func float64SliceValues(vals []float64) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_DoubleValue{
				DoubleValue: v,
			},
		}
	}
	return converted
}

func stringSliceValues(vals []string) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_StringValue{
				StringValue: v,
			},
		}
	}
	return converted
}

func InstrumentationScope(il instrumentation.Scope) *commonpb.InstrumentationScope {
	if il == (instrumentation.Scope{}) {
		return nil
	}
	return &commonpb.InstrumentationScope{
		Name:    il.Name,
		Version: il.Version,
	}
}
