package main

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace/examples"
	"fmt"
)

func main() {
	examples.StdoutTraceInit()
	//fmt.Println(examples.CheckAddressBefore())
	fmt.Println(examples.CheckAddress())
	examples.TraceProviderExit(context.Background())
}
