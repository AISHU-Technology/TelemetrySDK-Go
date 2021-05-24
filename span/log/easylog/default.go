package easylog

import (
	"os"
	"span/encoder"
	"span/field"
	"span/log"
	"span/open_standard"
	"span/runtime"
)

// return a Default SamplerLogger
func NewdefaultSamplerLogger() log.Logger {
	l := log.NewdefaultSamplerLogger()
	output := os.Stdin
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
	}
	writer.SetDefultResources()
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	return l
}
