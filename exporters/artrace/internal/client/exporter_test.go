package client

import (
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"sync"
	"testing"
)

func contextWithDone() context.Context {
	ctx := context.Background()
	done, cancel := context.WithCancel(ctx)
	cancel()
	return done
}

func TestExporter_ExportSpans(t *testing.T) {
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
				ross: nil,
			},
			false,
		}, {
			"StdoutClient_Exporter被停止",
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
				ross: nil,
			},
			false,
		}, {
			"HTTPClient_Exporter被停止",
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

func TestExporter_Shutdown(t *testing.T) {
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

var sClient = NewStdoutClient("")
var hClient = NewHTTPClient()
var sExporter = NewExporter(sClient)
var hExporter = NewExporter(hClient)

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
			args{sClient},
			sExporter,
		}, {
			"创建HTTPClient_Exporter",
			args{sClient},
			hExporter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got1, got2 := sExporter, hExporter; !reflect.DeepEqual(got1, tt.want) && !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("NewExporter() = %v, want %v", got1, tt.want)
			}
		})
	}
}
