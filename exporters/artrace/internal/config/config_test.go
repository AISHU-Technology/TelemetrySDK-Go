package config

import (
	"reflect"
	"testing"
)

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
			"没有配置项",
			args{opts: nil},
			Config{
				HTTPConfig:  DefaultHTTPConfig,
				RetryConfig: DefaultRetryConfig,
			},
		}, {
			"有配置项",
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
