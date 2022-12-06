package eventsdk

import (
	"reflect"
	"testing"
	"time"
)

var config = &eventProviderConfig{
	Exporters:     make(map[string]EventExporter),
	FlushInternal: 5 * time.Second,
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

func TestWithFlushInternal(t *testing.T) {
	type args struct {
		flushInternal time.Duration
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{time.Minute},
			WithFlushInternal(time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithFlushInternal(tt.args.flushInternal); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("WithFlushInternal() = %v, want %v", got, tt.want)
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
