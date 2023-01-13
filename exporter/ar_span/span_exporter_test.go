package ar_span

import (
	"context"
	"reflect"
	"testing"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
)

func TestNewExporter(t *testing.T) {
	type args struct {
		c public.Client
	}
	tests := []struct {
		name string
		args args
		want *SpanExporter
	}{
		{
			"HTTPClient的LogExporter",
			args{c: public.NewHTTPClient()},
			NewExporter(public.NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestLogExporterExportLogs(t *testing.T) {
	type fields struct {
		Exporter *public.Exporter
	}
	type args struct {
		ctx context.Context
		log []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"发送空数据",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: context.Background(),
				log: nil,
			},
			false,
		},
		{
			"发送Log",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: context.Background(),
				log: []byte("test"),
			},
			false,
		},
		{
			"已关闭的Exporter不能发Log",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: contextWithDone(),
				log: []byte("test"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SpanExporter{
				Exporter: tt.fields.Exporter,
			}
			if err := e.ExportLogs(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("ExportLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
