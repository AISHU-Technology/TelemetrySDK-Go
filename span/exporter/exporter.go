package exporter

import "context"

type LogExporter interface {
	ExportLogs(ctx context.Context, p []byte) error
	Name() string
	Shutdown(ctx context.Context) error
}
