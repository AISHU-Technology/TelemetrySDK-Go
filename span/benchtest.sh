#!/bin/bash

go test -benchmem -run=^$  -bench . gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/benchmarks -v
