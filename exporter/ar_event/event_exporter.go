package ar_event

import (
	"bytes"
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/eventsdk"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/public"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/resource"
	"encoding/json"
)

// 跨包实现接口占位用。
var _ eventsdk.EventExporter = (*EventExporter)(nil)

// EventExporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type EventExporter struct {
	*public.Exporter
}

// ExportEvents 批量发送 AnyRobotEvents 到AnyRobot Feed Ingester的 Event 数据接收器。
func (e *EventExporter) ExportEvents(ctx context.Context, events []eventsdk.Event) error {
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

// NewExporter 创建已启动的 EventExporter 。
func NewExporter(c public.Client) *EventExporter {
	return &EventExporter{
		public.NewExporter(c),
	}
}

// EventResource 传入 Event 的默认resource。
func EventResource() eventsdk.EventProviderOption {
	return resource.EventResource()
}
