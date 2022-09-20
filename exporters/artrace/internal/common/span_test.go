package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"testing"
)

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
			"TestAnyRobotSpanFromReadOnlySpan",
			args{},
			&AnyRobotSpan{},
		},
		{
			"TestAnyRobotSpanFromReadOnlySpan",
			args{nil},
			&AnyRobotSpan{},
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
			"TestAnyRobotSpansFromReadOnlySpans",
			args{},
			[]*AnyRobotSpan{},
		},
		{
			"TestAnyRobotSpansFromReadOnlySpans",
			args{nil},
			[]*AnyRobotSpan{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotSpansFromReadOnlySpans(tt.args.ross); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotSpansFromReadOnlySpans() = %v, want %v", got, tt.want)
			}
		})
	}
}