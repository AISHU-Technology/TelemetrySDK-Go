package artrace

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"go.opentelemetry.io/otel/sdk/resource"
	"reflect"
	"testing"
	"time"
)

var r1 = GetResource("name", "version")

func TestGetResource(t *testing.T) {
	type args struct {
		serviceName    string
		serviceVersion string
	}
	tests := []struct {
		name string
		args args
		want *resource.Resource
	}{
		{
			"填写服务名和版本后，获取资源信息",
			args{
				serviceName:    "name",
				serviceVersion: "version",
			},
			r1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r1; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

var eh = NewExporter(NewHTTPClient())

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
			"创建默认的HTTP发送器",
			args{NewHTTPClient()},
			eh,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eh; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

var ch = NewHTTPClient()

func TestNewHTTPClient(t *testing.T) {
	type args struct {
		opts []config.HTTPOption
	}
	tests := []struct {
		name string
		args args
		want client.Client
	}{
		{
			"创建HTTP Client",
			args{nil},
			ch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ch; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

var cs = NewStdoutClient("")

func TestNewStdoutClient(t *testing.T) {
	type args struct {
		stdoutPath string
	}
	tests := []struct {
		name string
		args args
		want client.Client
	}{
		{
			"创建Stdout Client",
			args{""},
			cs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cs; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

var optionWithAnyRobotURL = WithAnyRobotURL("https://www.baidu.com/")

func TestWithAnyRobotURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want *config.HTTPOption
	}{
		{
			"设置正确的上报地址",
			args{"https://www.baidu.com/"},
			&optionWithAnyRobotURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := &optionWithAnyRobotURL; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

var optionWithCompression = WithCompression(0)

func TestWithCompression(t *testing.T) {
	type args struct {
		compression int
	}
	tests := []struct {
		name string
		args args
		want *config.HTTPOption
	}{
		{
			"设置压缩方式",
			args{
				0,
			},
			&optionWithCompression,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := &optionWithCompression; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompression() = %v, want %v", got, tt.want)
			}
		})
	}
}

var optionWithHeader = WithHeader(nil)

func TestWithHeader(t *testing.T) {
	type args struct {
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want *config.HTTPOption
	}{
		{
			"设置自定义请求头信息",
			args{
				nil,
			},
			&optionWithHeader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := &optionWithHeader; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

var optionWithRetry = WithRetry(false, 0, 0, 0)

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
		want *config.HTTPOption
	}{
		{
			"关闭重发机制",
			args{
				enabled:        false,
				internal:       0,
				maxInterval:    0,
				maxElapsedTime: 0,
			},
			&optionWithRetry,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := &optionWithRetry; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

var optionWithTimeout = WithTimeout(0)

func TestWithTimeout(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want *config.HTTPOption
	}{
		{
			"设置超时时间",
			args{0},
			&optionWithTimeout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := &optionWithTimeout; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}
