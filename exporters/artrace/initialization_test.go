package artrace

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"reflect"
	"testing"
	"time"
)

func TestNewExporter(t *testing.T) {
	type args struct {
		c client.Client
	}
	tests := []struct {
		name string
		args args
		want *client.Exporter
	}{
		{
			name: "TestNewExporter_1",
			args: args{c: NewStdoutClient()},
			want: NewExporter(NewStdoutClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStdoutClient(t *testing.T) {
	tests := []struct {
		name string
		want client.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdoutClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNewHTTPClient(t *testing.T) {
	type args struct {
		opts []config.HTTPOption
	}
	tests := []struct {
		name string
		args args
		want client.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestWithAnyRobotURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want config.HTTPOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAnyRobotURL(tt.args.URL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	type args struct {
		compression int
	}
	tests := []struct {
		name string
		args args
		want config.HTTPOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompression(tt.args.compression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompression() = %v, want %v", got, tt.want)
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
		want config.HTTPOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.duration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
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
		want config.HTTPOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeader(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	type args struct {
		enabled        bool
		internal       time.Duration
		maxInterval    time.Duration
		maxElapsedTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want config.HTTPOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetry(tt.args.enabled, tt.args.internal, tt.args.maxInterval, tt.args.maxElapsedTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}
