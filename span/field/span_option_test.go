package field

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestWithAttribute(t *testing.T) {
	attr := NewAttribute("test", StringField("test_demo"))
	opt := WithAttribute(attr)
	info := &logSpanV1{}
	opt(info)
	assert.NotEqual(t, attr, info.attributes)
	record := MallocStructField(2)
	record.Set(attr.Type, attr.Message)
	//record.Set("Type", StringField(attr.Type))
	assert.Equal(t, record, info.attributes)

	opt1 := WithAttribute(nil)
	info1 := &logSpanV1{}
	opt1(info1)
	assert.Equal(t, MallocStructField(0), info1.GetAttributes())
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	opt := WithContext(ctx)
	info := &logSpanV1{}
	opt(info)
	assert.Equal(t, ctx, info.ctx)

	tp1 := tracesdk.NewTracerProvider()
	tr1 := tp1.Tracer("123")
	ctx1, span := tr1.Start(context.Background(), "fdsaf")
	defer tp1.Shutdown(nil)
	defer span.End()
	assert.NotEqual(t, info.ctx, ctx1)

}
