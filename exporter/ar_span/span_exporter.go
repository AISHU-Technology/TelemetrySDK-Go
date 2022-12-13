package ar_span

import (
	"context"
	"io"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
)

var _ io.Writer = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 log 数据接收器。
type Exporter struct {
	*public.Exporter
}

// ExportEvents 批量发送 AnyRobotEvents 到AnyRobot Feed Ingester的 Event 数据接收器。
// func (e *Exporter) ExportEvents(ctx context.Context, events []eventsdk.Event) error {
// 	if len(events) == 0 {
// 		return nil
// 	}
// 	file := bytes.NewBuffer([]byte{})
// 	encoder := json.NewEncoder(file)
// 	encoder.SetEscapeHTML(false)
// 	encoder.SetIndent("", "\t")
// 	if err := encoder.Encode(events); err != nil {
// 		return err
// 	}
// 	return e.ExportData(ctx, file.Bytes())
// }

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}

// Write 批量发送 log 到AnyRobot Feed Ingester的 log 数据接收器。
func (e *Exporter) Write(p []byte) (n int, err error) {
	ctx := context.Background()
	exportDataErr := e.ExportData(ctx, p)
	if err != nil {
		return 0, exportDataErr
	} else {
		return len(p), nil
	}
}
