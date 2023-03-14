package exporter

import "context"

type LogExporter interface {
	Name() string
	ExportLogs(ctx context.Context, p []byte) error
	Shutdown(ctx context.Context) error
}
