package eventsdk

import (
	"reflect"
	"testing"
	"time"
)

var config = &eventProviderConfig{
	Exporters:     make(map[string]EventExporter),
	FlashInternal: 5 * time.Second,
	MaxEvent:      9,
}

func TestWithExporters(t *testing.T) {
	type args struct {
		exporters []EventExporter
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{nil},
			WithExporters(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithExporters(tt.args.exporters...); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("WithExporters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFlashInternal(t *testing.T) {
	type args struct {
		flashInternal time.Duration
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{time.Minute},
			WithFlashInternal(time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithFlashInternal(tt.args.flashInternal); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("WithFlashInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxEvent(t *testing.T) {
	type args struct {
		maxEvent int
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{99},
			WithMaxEvent(99),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMaxEvent(tt.args.maxEvent); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("WithMaxEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServiceInfo(t *testing.T) {
	type args struct {
		ServiceName     string
		ServiceVersion  string
		ServiceInstance string
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{
				ServiceName:     "",
				ServiceVersion:  "",
				ServiceInstance: "",
			},
			WithServiceInfo("", "", ""),
		}, {
			"",
			args{
				ServiceName:     "123",
				ServiceVersion:  "456",
				ServiceInstance: "789",
			},
			WithServiceInfo("123", "456", "789"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithServiceInfo(tt.args.ServiceName, tt.args.ServiceVersion, tt.args.ServiceInstance); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("WithServiceInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventProviderOptionFuncApply(t *testing.T) {
	type args struct {
		cfg *eventProviderConfig
	}
	tests := []struct {
		name string
		o    eventProviderOptionFunc
		args args
		want *eventProviderConfig
	}{
		{
			"",
			eventProviderOptionFunc(func(providerConfig *eventProviderConfig) *eventProviderConfig {
				return providerConfig
			}),
			args{config},
			config,
		}, {
			"",
			eventProviderOptionFunc(func(providerConfig *eventProviderConfig) *eventProviderConfig {
				providerConfig.MaxEvent = 12
				return providerConfig
			}),
			args{WithMaxEvent(12).apply(config)},
			config,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.apply(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
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
				FlashInternal: 5 * time.Second,
				MaxEvent:      9,
			},
		}, {
			"",
			args{[]EventProviderOption{WithMaxEvent(19)}},
			WithMaxEvent(19).apply(&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlashInternal: 5 * time.Second,
				MaxEvent:      9,
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
