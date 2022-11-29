package custom_errors

// custom_errors 定义错误码和错误描述。
const (
	ModuleName = "TelemetrySDK-EventProvider(Go).Error: "

	Event_InvalidKey      = ModuleName + "Attribute设置了无意义的空键或与默认值冲突"
	Event_InvalidJSON     = ModuleName + "传入了非法的JSON，应该传入[]Event类型"
	Event_AlreadyShutdown = ModuleName + "已经关闭了Event Exporter，不能再发送"
)
