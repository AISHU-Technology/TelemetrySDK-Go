package artrace

import (
	"context"
	"testing"
	"time"
)

func TestInstallExportPipeline(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    func(context.Context) error
		wantErr bool
	}{
		{
			"TestInstallExportPipeline",
			args{
				context.Background(),
			},
			func(ctx context.Context) error {
				return nil
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InstallExportPipeline()
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallExportPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("InstallExportPipeline() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestSetAnyRobotURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestSetAnyRobotURL",
			args{
				"http://10.4.130.68:880/api/feed_ingester/v1/jobs/traceTest/events",
			},
			false,
		},
		{
			"TestSetAnyRobotURL",
			args{
				"://10.",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetAnyRobotURL(tt.args.URL); (err != nil) != tt.wantErr {
				t.Errorf("SetAnyRobotURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCompression(t *testing.T) {
	type args struct {
		Compression int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestSetCompression",
			args{
				0,
			},
			false,
		},
		{
			"TestSetCompression",
			args{
				4,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetCompression(tt.args.Compression); (err != nil) != tt.wantErr {
				t.Errorf("SetCompression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetInstrumentation(t *testing.T) {
	type args struct {
		InstrumentationName    string
		InstrumentationVersion string
		InstrumentationURL     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestSetInstrumentation",
			args{
				"go.opentelemetry.io/otel",
				"v1.9.0",
				"https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.9.0",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetInstrumentation(tt.args.InstrumentationName, tt.args.InstrumentationVersion, tt.args.InstrumentationURL); (err != nil) != tt.wantErr {
				t.Errorf("SetInstrumentation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetRetry(t *testing.T) {
	type args struct {
		internal       time.Duration
		maxInterval    time.Duration
		maxElapsedTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestSetRetry",
			args{
				time.Second,
				time.Second,
				time.Second,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetRetry(tt.args.internal, tt.args.maxInterval, tt.args.maxElapsedTime); (err != nil) != tt.wantErr {
				t.Errorf("SetRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetServiceResource(t *testing.T) {
	type args struct {
		url     string
		name    string
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestSetServiceResource",
			args{
				"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/artrace",
				"AnyRobotTrace-example",
				"2.2.0",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetServiceResource(tt.args.url, tt.args.name, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("SetServiceResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
