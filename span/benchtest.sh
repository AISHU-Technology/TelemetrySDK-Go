#!/bin/bash

go test -benchmem -run=^$  -bench . devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/benchmarks -v
