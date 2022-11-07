[English version](README_en.md)

![LOGO](docs/images/TelemetrySDK.png)

# [TelemetrySDK-Go Trace](https://devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go?version=GB2.2.0&path=/exporters/artrace/README.md&_a=preview)

`TelemetrySDK-Go`是 [OpenTelemetry](https://opentelemetry.io/) 的 [Go](https://golang.org/)
语言版本实现。本项目提供了一系列接口帮助开发者完成代码埋点过程，旨在提高用户业务的可观测性能力。

`Trace`是可观测性建设中的分布式链路追踪环节。文档旨在介绍分布式链路追踪的相关概念并指导如何使用`TelemetrySDK-Go Trace`
在Go语言编写的业务代码中埋点并上报链路数据到AnyRobot来建设可观测性能力。

### 兼容性

TelemetrySDK-Go 要求Go版本不低于`1.18`。

### [术语表](docs/cn/glossary.md)

### [开发指南](docs/cn/dev_guide.md)

### [接口文档](docs/cn/interface.md)

### [改造示例](examples/oneservice.go)

### 最佳实践