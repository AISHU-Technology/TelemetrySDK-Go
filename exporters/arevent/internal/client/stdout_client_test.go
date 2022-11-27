package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"reflect"
	"testing"
)

func TestNewStdoutClient(t *testing.T) {
	type args struct {
		stdoutPath string
	}
	tests := []struct {
		name string
		args args
		want EventClient
	}{
		{
			"",
			args{stdoutPath: ""},
			NewStdoutClient(""),
		}, {
			"",
			args{stdoutPath: "./simple.rst"},
			NewStdoutClient("./simple.rst"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdoutClient(tt.args.stdoutPath); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdoutClientPath(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				filepath: "/path",
				stopCh:   nil,
			},
			"/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &stdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			"",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{context.Background()},
			false,
		}, {
			"",
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

func TestStdoutClientUploadEvents(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
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
				filepath: "./event.txt",
				stopCh:   make(chan struct{}),
			},
			args{
				context.Background(),
				make([]eventsdk.Event, 1),
			},
			false,
		}, {
			"",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{
				contextWithDone(),
				[]eventsdk.Event{},
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
			if err := d.UploadEvents(tt.args.ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("UploadEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func contextWithCancelFunc() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	return ctx
}

func cancelFuncWithContext() context.CancelFunc {
	_, cancel := context.WithCancel(context.Background())
	return cancel
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
				stopCh:   channelWithStop(),
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
				t.Errorf("contextWithStop() got1 = %v, want %v", &got1, &tt.want1)
			}
		})
	}
}
