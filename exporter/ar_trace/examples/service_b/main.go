package main

import "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace/examples"

func main() {
	examples.StdoutTraceInit()
	//examples.ServerBefore()
	examples.Server()
}
