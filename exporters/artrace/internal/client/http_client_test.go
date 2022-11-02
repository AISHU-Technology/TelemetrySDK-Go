package client

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/common"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"github.com/agiledragon/gomonkey/v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func configWithStatusCode(statusCode string) config.HTTPConfig {
	cfg := config.DefaultHTTPConfig
	cfg.Path = statusCode
	cfg.Insecure = false
	cfg.Headers = map[string]string{"self": "defined"}
	cfg.Compression = 0
	return cfg
}

func TestHttpClientStop(t *testing.T) {
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"停止运行中的HttpClient",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{context.Background()},
			false,
		}, {
			"停止被context关闭的HttpClient",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpClientUploadTraces(t *testing.T) {
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		ctx           context.Context
		AnyRobotSpans []*common.AnyRobotSpan
	}

	sth := gomonkey.ApplyFunc(send, func(d *HttpClient, req *http.Request) (*http.Response, error) {
		if d.cfg.Path == "500" {
			code := 500
			return &http.Response{StatusCode: code,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte{})),
				Header: map[string][]string{"Retry-After": {"12"}}}, config.RetryableError{Throttle: 12}
		}
		if len(d.cfg.Path) == 3 {
			code, _ := strconv.Atoi(d.cfg.Path)
			return &http.Response{StatusCode: code,
				Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}, nil
		}
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}, nil
	})
	defer sth.Reset()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"发送非空Trace",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			false,
		}, {
			"返回错误码200",
			fields{
				configWithStatusCode("200"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			false,
		}, {
			"返回错误码400",
			fields{
				configWithStatusCode("400"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"返回错误码404",
			fields{
				configWithStatusCode("404"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"返回错误码413",
			fields{
				configWithStatusCode("413"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"返回错误码429",
			fields{
				configWithStatusCode("429"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"返回错误码500",
			fields{
				configWithStatusCode("500"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"返回错误码511",
			fields{
				configWithStatusCode("511"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			true,
		}, {
			"已关闭StdoutClient，不发送Trace",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				ctx:           contextWithDone(),
				AnyRobotSpans: nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if err := d.UploadTraces(tt.args.ctx, tt.args.AnyRobotSpans); (err != nil) != tt.wantErr {
				t.Errorf("UploadTraces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpClientContextWithStop(t *testing.T) {
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
		want1  context.CancelFunc
	}{
		{
			"正常返回等待执行的CancelFunc",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		}, {
			"被context关闭的不执行CancelFunc",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{ctx: contextWithDone()},
			contextWithDone(),
			cancelFuncWithContext(),
		}, {
			"被stopCh关闭立即执行CancelFunc",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				channelWithClosed(),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			got, got1 := d.contextWithStop(tt.args.ctx)
			if !reflect.DeepEqual(got.Err(), tt.want.Err()) {
				t.Errorf("contextWithStop() got = %v, want %v", got, tt.want)
			}
			if got1 == nil || tt.want1 == nil {
				t.Errorf("contextWithStop() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
