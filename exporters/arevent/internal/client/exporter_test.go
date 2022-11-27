package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"reflect"
	"sync"
	"testing"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func channelWithStop() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func TestExporterExportEvents(t *testing.T) {
	type fields struct {
		name     string
		client   EventClient
		stopCh   chan struct{}
		stopOnce *sync.Once
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
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		}, {
			"",
			fields{
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			false,
		}, {
			"",
			fields{
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				contextWithDone(),
				nil,
			},
			true,
		}, {
			"",
			fields{
				"",
				NewStdoutClient(""),
				channelWithStop(),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		}, {
			"",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		}, {
			"",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			true,
		}, {
			"",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				contextWithDone(),
				nil,
			},
			true,
		}, {
			"",
			fields{
				"",
				NewHTTPClient(),
				channelWithStop(),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.ExportEvents(tt.args.ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("ExportEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExporterName(t *testing.T) {
	type fields struct {
		name     string
		client   EventClient
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				name:     "NAME",
				client:   nil,
				stopCh:   nil,
				stopOnce: &sync.Once{},
			},
			"NAME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if got := e.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	type fields struct {
		name     string
		client   EventClient
		stopCh   chan struct{}
		stopOnce *sync.Once
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
				name:     "",
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				context.Background(),
			},
			false,
		}, {
			"",
			fields{
				name:     "",
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				contextWithDone(),
			},
			true,
		}, {
			"",
			fields{
				name:     "",
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				context.Background(),
			},
			false,
		}, {
			"",
			fields{
				name:     "",
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				contextWithDone(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
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
		client EventClient
	}
	tests := []struct {
		name string
		args args
		want *Exporter
	}{
		{
			"",
			args{NewStdoutClient("")},
			NewExporter(NewStdoutClient("")),
		}, {
			"",
			args{NewHTTPClient()},
			NewExporter(NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.client); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
