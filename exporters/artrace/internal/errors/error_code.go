package errors

// errors 定义错误码和错误描述。
const (
	ModuleName = "AnyRobotTraceExporter"

	AnyRobotTraceExporter_JobIdNotFound          = ModuleName + "接收器上报地址不正确"
	AnyRobotTraceExporter_PayloadTooLarge        = ModuleName + "Trace数据太大超过了5MB限制"
	AnyRobotTraceExporter_InvalidFormat          = ModuleName + "格式校验不通过"
	AnyRobotTraceExporter_InvalidURL             = ModuleName + "URL非法，请检查"
	AnyRobotTraceExporter_InvalidCompression     = ModuleName + "压缩方式不存在"
	AnyRobotTraceExporter_RetryTooLong           = ModuleName + "重发持续时间太长"
	AnyRobotTraceExporter_Unsent                 = ModuleName + "发送数据失败，检查日志"
	AnyRobotTraceExporter_ExceedRetryElapsedTime = ModuleName + "超过最大重发时间限制"
	AnyRobotTraceExporter_RetryFailure           = ModuleName + "Trace正在重发"
)
