package eventsdk

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultEventProviderConfig(t *testing.T) {
	tests := []struct {
		name string
		want *eventProviderConfig
	}{
		{
			"",
			&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: Internal,
				MaxEvent:      MaxEvent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultEventProviderConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultEventProviderConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventProviderConfig(t *testing.T) {
	type args struct {
		opts []EventProviderOption
	}
	tests := []struct {
		name string
		args args
		want *eventProviderConfig
	}{
		{
			"",
			args{nil},
			&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: 10 * time.Second,
				MaxEvent:      99,
			},
		}, {
			"",
			args{[]EventProviderOption{WithMaxEvent(19)}},
			WithMaxEvent(19).apply(&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: 10 * time.Second,
				MaxEvent:      99,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEventProviderConfig(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEventProviderConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
