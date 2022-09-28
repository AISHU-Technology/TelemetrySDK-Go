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
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"StdoutClient发送空trace",
			fields{
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
			"StdoutClient发送非空trace",
			fields{
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
			"2",
			fields{
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
			"3",
			fields{
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
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
				stopOnce: tt.fields.stopOnce,
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
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"1",
			fields{
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
			args{ctx: context.Background()},
			false,
		},
		{
			"2",
			fields{
				client:   NewStdoutClient(""),
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
				stopOnce: tt.fields.stopOnce,
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
