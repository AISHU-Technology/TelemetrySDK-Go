package client

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/common"
	"encoding/json"
	"io/ioutil"
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
		want Client
	}{
		{
			"创建StdoutClient",
			args{stdoutPath: ""},
			sClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sClient; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stdoutClient_Stop(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct {
		}
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
			"停止已关闭的StdoutClient",
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

func getAnyRobotSpans() []*common.AnyRobotSpan {
	var spans []*common.AnyRobotSpan
	bytes, _ := ioutil.ReadFile("./AnyRobotTrace.txt")
	_ = json.Unmarshal(bytes, &spans)
	return spans
}

func Test_stdoutClient_UploadTraces(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct {
		}
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
			"发送空Trace",
			fields{
				filepath: "./AnyRobotTrace.txt",
				stopCh:   make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: nil,
			},
			false,
		}, {
			"已关闭StdoutClient，不发送Trace",
			fields{
				filepath: "./AnyRobotTrace.txt",
				stopCh:   make(chan struct{}),
			},
			args{
				ctx:           contextWithDone(),
				AnyRobotSpans: getAnyRobotSpans(),
			},
			true,
		}, {
			"发送非空Trace",
			fields{
				filepath: "./AnyRobotTrace.txt",
				stopCh:   make(chan struct{}),
			},
			args{
				ctx:           context.Background(),
				AnyRobotSpans: getAnyRobotSpans(),
			},
			false,
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
