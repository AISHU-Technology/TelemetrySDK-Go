package ar_event

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"encoding/json"
)

var _ eventsdk.EventExporter = (*Exporter)(nil)

// Exporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type Exporter struct {
	*public.Exporter
}

// ExportEvents 批量发送 AnyRobotEvents 到AnyRobot Feed Ingester的 Event 数据接收器。
func (e *Exporter) ExportEvents(ctx context.Context, events []eventsdk.Event) error {
	if len(events) == 0 {
		return nil
	}
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(events); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// NewExporter 创建已启动的Exporter。
func NewExporter(c public.Client) *Exporter {
	return &Exporter{
		public.NewExporter(c),
	}
}
