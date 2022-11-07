package artrace

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/client"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"go.opentelemetry.io/otel/sdk/resource"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestGetResource(t *testing.T) {
	type args struct {
		serviceName       string
		serviceVersion    string
		serviceInstanceID string
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
			GetResource("name", "version", ""),
		}, {
			"填写服务名、版本、服务实例ID后，获取资源信息",
			args{
				serviceName:       "name",
				serviceVersion:    "version",
				serviceInstanceID: "edfd25ff-3c9c-b1a4-e660-bd826495ad35",
			},
			GetResource("name", "version", "edfd25ff-3c9c-b1a4-e660-bd826495ad35"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResource(tt.args.serviceName, tt.args.serviceVersion, tt.args.serviceInstanceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			"创建默认的Stdout发送器",
			args{NewStdoutClient("")},
			NewExporter(NewStdoutClient("")),
		}, {
			"创建默认的HTTP发送器",
			args{NewHTTPClient()},
			NewExporter(NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got.Shutdown(context.Background()), tt.want.Shutdown(context.Background())) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHTTPClient(t *testing.T) {
	type args struct {
		opts []config.Option
	}
	tests := []struct {
		name string
		args args
		want client.Client
	}{
		{
			"创建HTTP Client，不带Option",
			args{nil},
			NewHTTPClient(),
		}, {
			"创建HTTP Client，带上Option",
			args{[]config.Option{WithTimeout(0), WithCompression(1)}},
			NewHTTPClient(WithTimeout(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPClient(tt.args.opts...); !reflect.DeepEqual(got.Stop(context.Background()), tt.want.Stop(context.Background())) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			"创建Stdout Client，不指定输出文件地址",
			args{""},
			NewStdoutClient(""),
		}, {
			"创建Stdout Client，指定输出文件地址",
			args{"../AnyRobotTrace.txt"},
			NewStdoutClient("../AnyRobotTrace.txt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdoutClient(tt.args.stdoutPath); !reflect.DeepEqual(got.Stop(context.Background()), tt.want.Stop(context.Background())) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAnyRobotURL(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置正确的上报地址",
			args{"https://www.baidu.com/"},
			WithAnyRobotURL("https://www.baidu.com/"),
		}, {
			"设置错误的上报地址",
			args{"1:2:3:4:5"},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAnyRobotURL(tt.args.URL); !reflect.DeepEqual(got.Fn(config.DefaultConfig), tt.want.Fn(config.DefaultConfig)) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	fmt.Println()
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		compression int
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置压缩方式为JSON",
			args{
				0,
			},
			WithCompression(0),
		}, {
			"设置压缩方式为GZIP",
			args{
				1,
			},
			WithCompression(1),
		}, {
			"设置压缩方式为非法",
			args{
				2,
			},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompression(tt.args.compression); !reflect.DeepEqual(got.Fn(config.DefaultConfig), tt.want.Fn(config.DefaultConfig)) {
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
		want config.Option
	}{
		{
			"设置空请求头信息",
			args{
				nil,
			},
			WithHeader(nil),
		}, {
			"设置自定义请求头信息",
			args{
				map[string]string{"key": "value"},
			},
			WithHeader(map[string]string{"key": "value"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeader(tt.args.headers); !reflect.DeepEqual(got.Fn(config.DefaultConfig), tt.want.Fn(config.DefaultConfig)) {
				t.Errorf("WithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		enabled        bool
		internal       time.Duration
		maxInterval    time.Duration
		maxElapsedTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"关闭重发机制",
			args{
				enabled:        false,
				internal:       1,
				maxInterval:    2,
				maxElapsedTime: 3,
			},
			WithRetry(false, 1, 2, 3),
		}, {
			"开启重发机制",
			args{
				enabled:        true,
				internal:       4,
				maxInterval:    5,
				maxElapsedTime: 6,
			},
			WithRetry(true, 4, 5, 6),
		}, {
			"重发机制超时",
			args{
				enabled:        true,
				internal:       1 * time.Hour,
				maxInterval:    2 * time.Hour,
				maxElapsedTime: 3 * time.Hour,
			},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetry(tt.args.enabled, tt.args.internal, tt.args.maxInterval, tt.args.maxElapsedTime); !reflect.DeepEqual(got.Fn(config.DefaultConfig), tt.want.Fn(config.DefaultConfig)) {
				t.Errorf("WithRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置超时时间",
			args{2 * time.Second},
			WithTimeout(2 * time.Second),
		}, {
			"超时时间非法",
			args{61 * time.Second},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.duration); !reflect.DeepEqual(got.Fn(config.DefaultConfig), tt.want.Fn(config.DefaultConfig)) {
				t.Errorf("WithTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}
