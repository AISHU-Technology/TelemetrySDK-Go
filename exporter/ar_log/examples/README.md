# 代码示例说明

## 调用关系

single_service.go单独文件

## single_service.go

原始的包含同一个服务内提供加法、乘法处理函数。

## single_service_with_log.go

改造后的同一个服务内加法、乘法函数记录日志示例，目的是展示如何在代码中加入系统日志器、业务日志器的初始化配置。

# single_service运行过程

先阅读single_service.go代码，后运行single_service_test.go查看结果。
再阅读single_service_with_log.go代码，后运行single_service_with_log_test.go比较区别。查看AnyRobotTrace.json。

# 代码示例说明

## 调用关系

multi_service_a.go->multi_service_b.go->multi_service_c.go

## multi_service_a.go

原始的包含查询省份+城市的地址信息的函数。

## multi_service_b.go

使用了http框架的服务，运行在本地。分别提供查询省份接口和查询城市接口。

## multi_service_c.go

数据库查询服务，传入省份id或城市id，返回结果。需要连接本地数据库并填充简单示例。

| id  | address  |
|:---:|:--------:|
|  1  | ShangHai |
|  2  | BeiJing  |
|  3  | SiChuan  |
|  4  | ChengDu  |

## multi_service_a_with_log_and_trace.go

## multi_service_b_with_log_and_trace.go

改造后的模拟查询地址程序入口，目的是展示如何在代码中加入链路、日志的初始化配置，以及在日志中关联链路。

## multi_service_c_with_log_and_trace.go

改造后的数据库查询模拟。

# multi_service运行过程

1. 启动service_b/main.go，运行service_a/main.go，成功获取“Address : SiChuan Province ChengDu City”。
2. StdoutClient仅用于调试，正式使用应修改为HTTPClient。




