package config

import (
	"context"
	customErrors "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace/internal/errors"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"log"
	"reflect"
	"testing"
	"time"
)

func contextWithDone() context.Context {
	ctx := context.Background()
	done, cancel := context.WithCancel(ctx)
	cancel()
	return done
}

func TestRetryableErrorError(t *testing.T) {
	type fields struct {
		Throttle int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"生成重发报错",
			fields{0},
			"TelemetrySDK-Go-TraceExporter.Error: Trace正在重发",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := RetryableError{
				Throttle: tt.fields.Throttle,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 time.Duration
	}{
		{
			"error为空",
			args{nil},
			false,
			0,
		}, {
			"error为普通错误",
			args{errors.New("something")},
			false,
			0,
		}, {
			"error为空",
			args{RetryableError{15}},
			true,
			15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := evaluate(tt.args.err)
			if got != tt.want {
				t.Errorf("evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestWait(t *testing.T) {
	type args struct {
		ctx   context.Context
		delay time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"context正常且计时器到时间",
			args{
				ctx:   context.Background(),
				delay: 1 * time.Nanosecond,
			},
			false,
		}, {
			"context退出且未到时间",
			args{
				ctx:   contextWithDone(),
				delay: 100 * time.Nanosecond,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := wait(tt.args.ctx, tt.args.delay); (err != nil) != tt.wantErr {
				t.Errorf("wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	type args struct {
		rc RetryConfig
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"配置可选重发选项",
			args{DefaultRetryConfig},
			WithRetry(DefaultRetryConfig),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetry(tt.args.rc); !reflect.DeepEqual(got.applyOption(DefaultConfig), tt.want.applyOption(DefaultConfig)) {
				t.Errorf("WithRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

type key struct {
	name string
}

func deliver(ctx context.Context) error {
	switch ctx.Value(key{"StatusCode"}) {
	case 200:
		return nil
	case 413:
		return errors.New("413")
	case 500:
		return RetryableError{Throttle: 10}
	case 503:
		return RetryableError{Throttle: 1}
	default:
		return nil
	}
}

func result(ctx context.Context) error {
	switch ctx.Value(key{"StatusCode"}) {
	case 200:
		return nil
	case 413:
		return errors.New("413")
	case 500:
		return errors.New(customErrors.AnyRobotTraceExporter_ExceedRetryElapsedTime)
	case 503:
		return RetryableError{}
	default:
		return nil
	}
}

func TestRetryConfigRetryFunc(t *testing.T) {
	sth := gomonkey.ApplyFunc(wait, func(ctx context.Context, delay time.Duration) error {
		return RetryableError{}
	})
	defer sth.Reset()

	other := gomonkey.ApplyFunc(log.Println, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer other.Reset()

	type fields struct {
		Enabled         bool
		InitialInterval time.Duration
		MaxInterval     time.Duration
		MaxElapsedTime  time.Duration
		ctx             context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   RetryFunc
	}{
		{
			"不开启重发",
			fields{
				Enabled:         false,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.Background(),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		}, {
			"启用了重发，返回error=nil",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.WithValue(context.Background(), key{"StatusCode"}, 200),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		}, {
			"启用了重发，返回error=PayloadTooLarge",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.WithValue(context.Background(), key{"StatusCode"}, 413),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		}, {
			"启用了重发，返回error=RetryFailure，达到最长重发时间",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     1 * time.Second,
				MaxElapsedTime:  2 * time.Nanosecond,
				ctx:             context.WithValue(context.Background(), key{"StatusCode"}, 500),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		}, {
			"启用了重发，返回error=RetryFailure，可以继续重发",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     1 * time.Second,
				MaxElapsedTime:  2 * time.Second,
				ctx:             context.WithValue(context.Background(), key{"StatusCode"}, 503),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RetryConfig{
				Enabled:         tt.fields.Enabled,
				InitialInterval: tt.fields.InitialInterval,
				MaxInterval:     tt.fields.MaxInterval,
				MaxElapsedTime:  tt.fields.MaxElapsedTime,
			}
			if got := r.RetryFunc(); !reflect.DeepEqual(got(tt.fields.ctx, deliver), tt.want(tt.fields.ctx, result)) {
				t.Errorf("RetryFunc() = %v, want %v", got(tt.fields.ctx, deliver), tt.want(tt.fields.ctx, result))
			}
		})
	}
}
