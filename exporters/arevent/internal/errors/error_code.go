package errors

// errors 定义错误码和错误描述。
const (
	ModuleName                           = "TelemetrySDK-EventExporter(Go).Error: "
	EventExporter_JobIdNotFound          = ModuleName + "接收器上报地址不正确"
	EventExporter_PayloadTooLarge        = ModuleName + "Trace数据太大超过了5MB限制"
	EventExporter_InvalidFormat          = ModuleName + "格式校验不通过"
	EventExporter_InvalidURL             = ModuleName + "URL非法，请检查"
	EventExporter_InvalidCompression     = ModuleName + "压缩方式不存在"
	EventExporter_RetryTooLong           = ModuleName + "重发持续时间太长"
	EventExporter_Unsent                 = ModuleName + "发送数据失败，检查日志"
	EventExporter_ExceedRetryElapsedTime = ModuleName + "超过最大重发时间限制"
	EventExporter_DurationTooLong        = ModuleName + "超过最长连接时间限制"
	EventExporter_RetryFailure           = ModuleName + "Trace正在重发"
)
