package easylog

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/log"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
	"os"
)

// return a Default SamplerLogger
func NewDefaultSamplerLogger() log.Logger {
	l := log.NewDefaultSamplerLogger()
	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
	}
	writer.SetDefultResources()
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	return l
}
