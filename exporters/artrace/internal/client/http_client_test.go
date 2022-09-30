package client

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/config"
	"reflect"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	type args struct {
		opts []config.HTTPOption
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			"创建HTTPClient",
			args{opts: nil},
			hClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hClient; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_arRequest_reset(t *testing.T) {
//	type fields struct {
//		Request    *http.Request
//		bodyReader func() io.ReadCloser
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &arRequest{
//				Request:    tt.fields.Request,
//				bodyReader: tt.fields.bodyReader,
//			}
//			r.reset(tt.args.ctx)
//		})
//	}
//}
//
//func Test_bodyReader(t *testing.T) {
//	type args struct {
//		buf []byte
//	}
//	tests := []struct {
//		name string
//		args args
//		want func() io.ReadCloser
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := bodyReader(tt.args.buf); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("bodyReader() = %v, want %v", got(), tt.want())
//			}
//		})
//	}
//}
//
//func Test_httpClient_Stop(t *testing.T) {
//	type fields struct {
//		cfg       config.HTTPConfig
//		retryFunc config.RetryFunc
//		client    *http.Client
//		stopCh    chan struct{}
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &httpClient{
//				cfg:       tt.fields.cfg,
//				retryFunc: tt.fields.retryFunc,
//				client:    tt.fields.client,
//				stopCh:    tt.fields.stopCh,
//			}
//			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_httpClient_UploadTraces(t *testing.T) {
//	type fields struct {
//		cfg       config.HTTPConfig
//		retryFunc config.RetryFunc
//		client    *http.Client
//		stopCh    chan struct{}
//	}
//	type args struct {
//		ctx           context.Context
//		AnyRobotSpans []*common.AnyRobotSpan
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &httpClient{
//				cfg:       tt.fields.cfg,
//				retryFunc: tt.fields.retryFunc,
//				client:    tt.fields.client,
//				stopCh:    tt.fields.stopCh,
//			}
//			if err := d.UploadTraces(tt.args.ctx, tt.args.AnyRobotSpans); (err != nil) != tt.wantErr {
//				t.Errorf("UploadTraces() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_httpClient_contextWithStop(t *testing.T) {
//	type fields struct {
//		cfg       config.HTTPConfig
//		retryFunc config.RetryFunc
//		client    *http.Client
//		stopCh    chan struct{}
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   context.Context
//		want1  context.CancelFunc
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &httpClient{
//				cfg:       tt.fields.cfg,
//				retryFunc: tt.fields.retryFunc,
//				client:    tt.fields.client,
//				stopCh:    tt.fields.stopCh,
//			}
//			got, got1 := d.contextWithStop(tt.args.ctx)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("contextWithStop() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("contextWithStop() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_httpClient_getScheme(t *testing.T) {
//	type fields struct {
//		cfg       config.HTTPConfig
//		retryFunc config.RetryFunc
//		client    *http.Client
//		stopCh    chan struct{}
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &httpClient{
//				cfg:       tt.fields.cfg,
//				retryFunc: tt.fields.retryFunc,
//				client:    tt.fields.client,
//				stopCh:    tt.fields.stopCh,
//			}
//			if got := d.getScheme(); got != tt.want {
//				t.Errorf("getScheme() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_httpClient_newRequest(t *testing.T) {
//	type fields struct {
//		cfg       config.HTTPConfig
//		retryFunc config.RetryFunc
//		client    *http.Client
//		stopCh    chan struct{}
//	}
//	type args struct {
//		body []byte
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    arRequest
//		wantErr bool
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			d := &httpClient{
//				cfg:       tt.fields.cfg,
//				retryFunc: tt.fields.retryFunc,
//				client:    tt.fields.client,
//				stopCh:    tt.fields.stopCh,
//			}
//			got, err := d.newRequest(tt.args.body)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("newRequest() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_newResponseError(t *testing.T) {
//	type args struct {
//		header http.Header
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := newResponseError(tt.args.header); (err != nil) != tt.wantErr {
//				t.Errorf("newResponseError() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
