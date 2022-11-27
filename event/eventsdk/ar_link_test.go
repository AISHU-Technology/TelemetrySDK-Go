package eventsdk

import (
	"reflect"
	"testing"
)

func TestLinkGetSpanID(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.GetSpanID(); got != tt.want {
				t.Errorf("GetSpanID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkGetTraceID(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.GetTraceID(); got != tt.want {
				t.Errorf("GetTraceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkPrivate(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			l.private()
		})
	}
}

func TestNewLink(t *testing.T) {
	tests := []struct {
		name string
		want link
	}{
		{
			"",
			link{
				TraceID: "00000000000000000000000000000000",
				SpanID:  "0000000000000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLink(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
