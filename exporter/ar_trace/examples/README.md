# 代码示例说明

## 调用关系

single_service.go单独文件

## single_service.go

原始的包含同一个服务内提供加法、乘法处理函数。

## single_service_with_trace.go

改造后的同一个服务内加法、乘法函数在不是父子调用关系时，在同一条调用链上的代码埋点示例，目的是展示如何在代码中加入链路的初始化配置。

# single_service运行过程

先阅读single_service.go代码，后运行single_service_test.go查看结果。
再阅读single_service_with_trace.go代码，后运行single_service_with_trace_test.go比较区别。查看AnyRobotTrace.json。

# 代码示例说明

## 调用关系

multi_service_a_http.go->multi_service_b_http.go->multi_service_c_mysql.go

## multi_service_a_http.go | multi_service_d_gin.go

原始的包含查询省份+城市的地址信息的函数。

## multi_service_b_http.go | multi_service_e_gin.go

使用了http框架的服务，运行在本地。分别提供查询省份接口和查询城市接口。

## multi_service_c_mysql.go

数据库查询服务，传入省份id或城市id，返回结果。需要连接本地数据库并填充简单示例。

| id  | address  |
|:---:|:--------:|
|  1  | ShangHai |
|  2  | BeiJing  |
|  3  | SiChuan  |
|  4  | ChengDu  |

## multi_service_a_with_trace_http.go | multi_service_d_with_trace_gin.go

改造后的模拟查询地址程序入口，目的是展示父子关系调用，如何在代码中加入链路的初始化配置以及在Go服务调用链中跨服务传播链路上下文信息。

## multi_service_b_with_trace_http.go | multi_service_e_with_trace_gin.go

改造后的http框架，加入了传播链路上下文信息的插件：otelhttp、
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))。

## multi_service_c_with_trace_mysql.go

改造后的数据库查询模拟。

# multi_service运行过程

1. 启动service_b/main.go，运行service_a/main.go，成功获取“Address : SiChuan Province ChengDu City”。
2. StdoutClient仅用于调试，正式使用应修改为HTTPClient。
3. 启动service_e/main.go，运行service_d/main.go，成功获取“Address : SiChuan Province ChengDu City”。
4. StdoutClient仅用于调试，正式使用应修改为HTTPClient。




