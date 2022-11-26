package errors

// errors 定义错误码和错误描述。
const (
	ModuleName = "TelemetrySDK-Go-EventExporter.Error: "

	AnyRobotEventExporter_JobIdNotFound          = ModuleName + "接收器上报地址不正确"
	AnyRobotEventExporter_PayloadTooLarge        = ModuleName + "Trace数据太大超过了5MB限制"
	AnyRobotEventExporter_InvalidFormat          = ModuleName + "格式校验不通过"
	AnyRobotEventExporter_InvalidURL             = ModuleName + "URL非法，请检查"
	AnyRobotEventExporter_InvalidCompression     = ModuleName + "压缩方式不存在"
	AnyRobotEventExporter_RetryTooLong           = ModuleName + "重发持续时间太长"
	AnyRobotEventExporter_Unsent                 = ModuleName + "发送数据失败，检查日志"
	AnyRobotEventExporter_ExceedRetryElapsedTime = ModuleName + "超过最大重发时间限制"
	AnyRobotEventExporter_DurationTooLong        = ModuleName + "超过最长连接时间限制"
	AnyRobotEventExporter_RetryFailure           = ModuleName + "Trace正在重发"
	AnyRobotEventExporter_InvalidKey             = ModuleName + "Attribute设置了无意义的空键或与默认值冲突"
	AnyRobotEventExporter_InvalidJSON            = ModuleName + "传入了非法的JSON，应该传入[]model.AREvent类型"
)
