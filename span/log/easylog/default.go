package easylog

import (
	"os"

	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/encoder"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/field"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/log"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/open_standard"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo/span/runtime"
)

// return a Default SamplerLogger
func NewDefaultSamplerLogger() log.Logger {
	l := log.NewDefaultSamplerLogger()
	output := os.Stdout
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
	}
	writer.SetDefaultResources()
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	return l
}
