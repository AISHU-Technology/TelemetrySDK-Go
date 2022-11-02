package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/common"
	"reflect"
	"testing"
)

func TestStdoutClientStop(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
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
			"停止运行中的StdoutClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{context.Background()},
			false,
		}, {
			"停止被context关闭的StdoutClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &stdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStdoutClientUploadTraces(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
	}
	type args struct {
		ctx           context.Context
		AnyRobotSpans []*common.AnyRobotSpan
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"发送非空Trace",
			fields{
				filepath: "./trace.txt",
				stopCh:   make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: []*common.AnyRobotSpan{{}, {}},
			},
			false,
		}, {
			"已关闭StdoutClient，不发送Trace",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
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
			d := &stdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.UploadTraces(tt.args.ctx, tt.args.AnyRobotSpans); (err != nil) != tt.wantErr {
				t.Errorf("UploadTraces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStdoutClientContextWithStop(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
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
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		}, {
			"被context关闭的不执行CancelFunc",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{ctx: contextWithDone()},
			contextWithDone(),
			cancelFuncWithContext(),
		}, {
			"被stopCh关闭立即执行CancelFunc",
			fields{
				filepath: "",
				stopCh:   channelWithClosed(),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &stdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
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

func TestNewStdoutClient(t *testing.T) {
	type args struct {
		stdoutPath string
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			"创建StdoutClient",
			args{stdoutPath: ""},
			NewStdoutClient(""),
		}, {
			"创建StdoutClient",
			args{stdoutPath: "./simple.rst"},
			NewStdoutClient("./simple.rst"),
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
