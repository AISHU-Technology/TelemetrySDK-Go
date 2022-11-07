package client

import (
	"context"
	"fmt"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"sync"
	"testing"
)

type MySpan struct {
	ros sdktrace.ReadOnlySpan
}

func MockReadOnlySpan() sdktrace.ReadOnlySpan {
	var mySpan = MySpan{}
	return mySpan.ros
}
func MockReadOnlySpans() []sdktrace.ReadOnlySpan {
	return []sdktrace.ReadOnlySpan{MockReadOnlySpan(), MockReadOnlySpan()}
}

func contextWithDone() context.Context {
	ctx := context.Background()
	done, cancel := context.WithCancel(ctx)
	cancel()
	return done
}

func channelWithClosed() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func cancelFuncWithContext() context.CancelFunc {
	_, cancel := context.WithCancel(context.Background())
	return cancel
}
func contextWithCancelFunc() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	defer fmt.Println(&cancel)
	return ctx
}

func TestExporterExportSpans(t *testing.T) {
	type fields struct {
		client   Client
		stopCh   chan struct{}
		stopOnce sync.Once
	}
	type args struct {
		ctx  context.Context
		ross []sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name    string
		fields  *fields
		args    args
		wantErr bool
	}{
		{
			"StdoutClient_Exporter发送空trace",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: nil,
			},
			false,
		}, {
			"StdoutClient_Exporter发送非空trace",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: MockReadOnlySpans(),
			},
			false,
		}, {
			"StdoutClient_Exporter被context停止",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  contextWithDone(),
				ross: nil,
			},
			true,
		}, {
			"StdoutClient_Exporter被stopCh停止",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   channelWithClosed(),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: nil,
			},
			false,
		}, {
			"HTTPClient_Exporter发送空trace",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: nil,
			},
			false,
		}, {
			"HTTPClient_Exporter发送非空trace",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: MockReadOnlySpans(),
			},
			true,
		}, {
			"HTTPClient_Exporter被context停止",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  contextWithDone(),
				ross: nil,
			},
			true,
		}, {
			"HTTPClient_Exporter被stopCh停止",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   channelWithClosed(),
				stopOnce: sync.Once{},
			},
			args{
				ctx:  context.Background(),
				ross: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.ExportSpans(tt.args.ctx, tt.args.ross); (err != nil) != tt.wantErr {
				t.Errorf("ExportSpans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	type fields struct {
		client   Client
		stopCh   chan struct{}
		stopOnce sync.Once
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  *fields
		args    args
		wantErr bool
	}{
		{
			"关闭运行中的StdoutClient_Exporter",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{ctx: context.Background()},
			false,
		}, {
			"关闭已经停止的StdoutClient_Exporter",
			&fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{ctx: contextWithDone()},
			true,
		}, {
			"关闭运行中的HTTPClient_Exporter",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{ctx: context.Background()},
			false,
		}, {
			"关闭已经停止的HTTPClient_Exporter",
			&fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{ctx: contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.Shutdown(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExporter(t *testing.T) {
	type args struct {
		client Client
	}
	tests := []struct {
		name string
		args args
		want *Exporter
	}{
		{
			"创建StdoutClient_Exporter",
			args{NewStdoutClient("")},
			NewExporter(NewStdoutClient("")),
		}, {
			"创建HTTPClient_Exporter",
			args{NewHTTPClient()},
			NewExporter(NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.client); !reflect.DeepEqual(got.client.Stop(context.Background()), tt.want.client.Stop(context.Background())) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
