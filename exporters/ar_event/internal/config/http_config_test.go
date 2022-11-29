package config

import (
	"reflect"
	"testing"
	"time"
)

func TestWithAnyRobotURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"",
			args{"http://www.rockyrori.cn/1234"},
			Option{Fn: func(cfg Config) Config {
				cfg.HTTPConfig.Insecure = true
				cfg.HTTPConfig.Endpoint = "www.rockyrori.cn"
				cfg.HTTPConfig.Path = "/1234"
				return cfg
			}},
		}, {
			"",
			args{"https://www.rockyrori.cn:80"},
			Option{Fn: func(cfg Config) Config {
				cfg.HTTPConfig.Insecure = false
				cfg.HTTPConfig.Endpoint = "www.rockyrori.cn:80"
				cfg.HTTPConfig.Path = ""
				return cfg
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAnyRobotURL(tt.args.URL); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	type args struct {
		compression Compression
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"",
			args{0},
			Option{Fn: func(cfg Config) Config {
				cfg.HTTPConfig.Compression = 0
				return cfg
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompression(tt.args.compression); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("WithCompression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHeader(t *testing.T) {
	type args struct {
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"",
			args{nil},
			Option{Fn: func(cfg Config) Config {
				cfg.HTTPConfig.Headers = nil
				return cfg
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeader(tt.args.headers); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("WithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"",
			args{123},
			Option{Fn: func(cfg Config) Config {
				cfg.HTTPConfig.Timeout = 123
				return cfg
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.duration); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}
