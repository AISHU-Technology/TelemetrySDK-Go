package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"testing"
)

type MySpan struct {
	ros sdktrace.ReadOnlySpan
}

func MockReadOnlySpan() sdktrace.ReadOnlySpan {
	var mySpan = MySpan{}
	return mySpan.ros
}
func MockReadOnlySpans() []sdktrace.ReadOnlySpan {
	return []sdktrace.ReadOnlySpan{MockReadOnlySpan(), MockReadOnlySpan()}
}

func TestAnyRobotSpanFromReadOnlySpan(t *testing.T) {
	type args struct {
		ros sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name string
		args args
		want *AnyRobotSpan
	}{
		{
			"转换空ReadOnlySpan",
			args{nil},
			AnyRobotSpanFromReadOnlySpan(nil),
		},
		{
			"转换非空ReadOnlySpan",
			args{MockReadOnlySpan()},
			AnyRobotSpanFromReadOnlySpan(MockReadOnlySpan()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotSpanFromReadOnlySpan(tt.args.ros); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotSpanFromReadOnlySpan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyRobotSpansFromReadOnlySpans(t *testing.T) {
	type args struct {
		ross []sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name string
		args args
		want []*AnyRobotSpan
	}{
		{
			"转换空ReadOnlySpan",
			args{nil},
			AnyRobotTraceFromReadOnlyTrace(nil),
		},
		{
			"转换非空ReadOnlySpans",
			args{MockReadOnlySpans()},
			AnyRobotTraceFromReadOnlyTrace(MockReadOnlySpans()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotTraceFromReadOnlyTrace(tt.args.ross); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotTraceFromReadOnlyTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}
