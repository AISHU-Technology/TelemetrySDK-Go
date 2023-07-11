package main

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_log/examples"
)

func main() {
	examples.StdoutTraceInit()
	examples.ConsoleLogInit()
	//examples.CheckAddressGinBefore()
	examples.CheckAddressGin()
	examples.TraceProviderExit(context.Background())
}
