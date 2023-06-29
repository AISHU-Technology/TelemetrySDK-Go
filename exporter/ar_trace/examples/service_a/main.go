package main

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace/examples"
	"fmt"
	"time"
)

func main() {
	examples.StdoutTraceInit()
	//fmt.Println(examples.CheckAddressBefore())
	fmt.Println(examples.CheckAddress())
	time.Sleep(5 * time.Second)
}
