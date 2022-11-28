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

func TestLinkValid(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				TraceID: "5bd0ecc145cc8639007721df27ecda50",
				SpanID:  "4cbf3e2c1e8517e7",
			},
			true,
		}, {
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
