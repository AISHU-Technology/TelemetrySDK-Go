package config

import (
	"reflect"
	"testing"
)

func TestEmptyOption(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		{
			"",
			Option{
				func(cfg Config) Config {
					return cfg
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyOption(); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("EmptyOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"",
			args{opts: nil},
			Config{
				HTTPConfig:  DefaultHTTPConfig,
				RetryConfig: DefaultRetryConfig,
			},
		}, {
			"",
			args{opts: []Option{EmptyOption()}},
			Config{
				HTTPConfig:  DefaultHTTPConfig,
				RetryConfig: DefaultRetryConfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionApplyOption(t *testing.T) {
	type fields struct {
		Fn func(Config) Config
	}
	type args struct {
		cfg Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Config
	}{
		{
			"",
			fields{EmptyOption().Fn},
			args{DefaultConfig},
			DefaultConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Option{
				Fn: tt.fields.Fn,
			}
			if got := h.applyOption(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOption(t *testing.T) {
	type args struct {
		fn func(cfg Config) Config
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"",
			args{},
			Option{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOption(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
