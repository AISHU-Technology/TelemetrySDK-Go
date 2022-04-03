#!/bin/bash

go test -benchmem -run=^$  -bench . devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/benchmarks -v
