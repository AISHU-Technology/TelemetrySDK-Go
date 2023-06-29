package main

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_log/examples"
	"fmt"
)

func main() {
	examples.StdoutTraceInit()
	examples.ConsoleLogInit()
	//fmt.Println(examples.CheckAddressBefore())
	fmt.Println(examples.CheckAddress())
	examples.TraceProviderExit(context.Background())
	examples.LoggerExit()
}
