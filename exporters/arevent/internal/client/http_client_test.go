package client

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/internal/config"
	"github.com/agiledragon/gomonkey/v2"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestHttpClientPath(t *testing.T) {
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				cfg:       config.HTTPConfig{},
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			"",
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
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
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
			"",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{context.Background()},
			false,
		}, {
			"",
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

func configWithStatusCode(statusCode string) config.HTTPConfig {
	cfg := config.DefaultHTTPConfig
	cfg.Path = statusCode
	cfg.Insecure = false
	cfg.Headers = map[string]string{"self": "defined"}
	cfg.Compression = 0
	return cfg
}

func TestHttpClientUploadEvents(t *testing.T) {
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

	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		ctx    context.Context
		events []eventsdk.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			false,
		}, {
			"",
			fields{
				configWithStatusCode("200"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			false,
		}, {
			"",
			fields{
				configWithStatusCode("400"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				configWithStatusCode("404"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				configWithStatusCode("413"),
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				configWithStatusCode("429"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				configWithStatusCode("500"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				configWithStatusCode("511"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				contextWithDone(),
				make([]eventsdk.Event, 1),
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
			if err := d.UploadEvents(tt.args.ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("UploadEvents() error = %v, wantErr %v", err, tt.wantErr)
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
			"",
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
			"",
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
			"",
			fields{
				config.DefaultHTTPConfig,
				config.DefaultRetryConfig.RetryFunc(),
				http.DefaultClient,
				channelWithStop(),
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
				t.Errorf("contextWithStop() got1 = %v, want %v", &got1, &tt.want1)
			}
		})
	}
}

func TestHttpClientGetScheme(t *testing.T) {
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				cfg:       config.HTTPConfig{},
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			"https",
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
			if got := d.getScheme(); got != tt.want {
				t.Errorf("getScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpClientNewRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	type fields struct {
		cfg       config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    arRequest
		wantErr bool
	}{
		{
			"",
			fields{
				cfg:       config.HTTPConfig{},
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			args{
				[]byte{},
			},
			arRequest{
				Request:    r,
				bodyReader: bodyReader([]byte{}),
			},
			false,
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
			got, err := d.newRequest(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Method, tt.want.Method) {
				t.Errorf("newRequest() got = %v, want %v", got, tt.want)
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
		want EventClient
	}{
		{
			"",
			args{nil},
			NewHTTPClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPClient(tt.args.opts...); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARRequestReset(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	type fields struct {
		Request    *http.Request
		bodyReader func() io.ReadCloser
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				Request:    r,
				bodyReader: bodyReader([]byte{}),
			},
			args{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &arRequest{
				Request:    tt.fields.Request,
				bodyReader: tt.fields.bodyReader,
			}
			r.reset(tt.args.ctx)
		})
	}
}

func TestBodyReader(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want func() io.ReadCloser
	}{
		{
			"",
			args{
				[]byte{},
			},
			bodyReader([]byte{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bodyReader(tt.args.buf); !reflect.DeepEqual(got(), tt.want()) {
				t.Errorf("bodyReader() = %v, want %v", got(), tt.want())
			}
		})
	}
}

func TestNewResponseError(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := newResponseError(tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("newResponseError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSend(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, "http://127.0.0.1", nil)
	res, _ := http.DefaultClient.Do(req)
	type args struct {
		d   *HttpClient
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"",
			args{
				d: &HttpClient{
					cfg:       config.DefaultConfig.HTTPConfig,
					retryFunc: config.DefaultConfig.RetryConfig.RetryFunc(),
					stopCh:    make(chan struct{}),
					client:    http.DefaultClient,
				},
				req: req,
			},
			res,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := send(tt.args.d, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("send() got = %v, want %v", got, tt.want)
			}
		})
	}
}
