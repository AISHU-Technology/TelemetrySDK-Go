# 数据分层介绍

当前实现数据模型分层如下图

![](./images/data_model_v1-layer.png)



# AishuV0协议

序列化层为'\n'分隔的json文本。每个`InternalSpan`为一次日志输出，内部包含指标、日志事件和外部调用追踪。多个`InternalSpan`聚合可以用于进行内部调用追踪/线程追踪。各个字段如下表所述

## InternalSpan

| 字段名        | 类型                | 描述                                                         |
| ------------- | ------------------- | ------------------------------------------------------------ |
| TraceId       | string              | 分布式事务ID，一个trace含有多个span                          |
| SpanId        | string              | 子分布式事务ID，多个span组成trace。span间的父子关系和兄弟关系可以描述一次事务的内外调用关系 |
| StartTime     | unixTime            | span开始时间                                                 |
| EndTime       | unixTime            | span结束时间                                                 |
| Events        | Event(Array)        | 数据事件数组，一个event一般代表线程内的一次日志记录行为。events可以组成线程内的日志上下文 |
| metrics       | Metric(Array)       | 指标数组。代表一次线程执行的指标数据                         |
| externalSpans | ExternalSpan(Array) | 外部调用追踪数组，表示当前线程发起的多次外部调用             |

## Event

| 字段名  | 类型         | 描述                                                         |
| ------- | ------------ | ------------------------------------------------------------ |
| Level   | Int          | 日志级别，Trace->Fatal对应值为1->6                           |
| type    | string       | event类型，默认为空，非空值时由应用定义用于扩展数据类型      |
| message | string或任意 | type为空时为字符串，type非空时由上层应用定义，由对应定义者提供数据模型 |

## Metric

| 字段名     | 类型              | 描述                                                      |
| ---------- | ----------------- | --------------------------------------------------------- |
| Attributes | map[string]string | 指标属性，为map类型，key值类型为字符串、value同样为字符串 |
| Labels     | string(array)     | 标签数组，数组内元素为数组                                |
| 其他字段   | string            | key值为metric名称，value值为metric值                      |

## ExternalSpan

| 字段名           | 类型     | 描述                                                         |
| ---------------- | -------- | ------------------------------------------------------------ |
| TraceId          | string   | 分布式事务ID，一个trace含有多个span                          |
| SpanId           | string   | 子分布式事务ID，多个span组成trace。span间的父子关系和兄弟关系可以描述一次事务的内外调用关系 |
| StartTime        | unixTime | span开始时间                                                 |
| EndTime          | unixTime | span结束时间                                                 |
| InternalParentId | string   | 发起该外部调用的内部spanID                                   |
| Attributes       | 任意     | 填写span信息，未来根据需要追踪的协议再次封装对应的链路追踪信息 |



### 理想数据样例

以下数据为无外部干扰的理想数据，数据内容为一个trace的两个InternalSpan，每个InternalSpan中含有一些ExternalSpan外部调用、metric以及应用业务事件日志。为方便介绍，每个日志后会额外打印一行空行，真实数据不存在该空行，即日志内容仅由一个换行符分隔

```text
{"TraceId": "1743fb330000100000000000000000300000000000000000","SpanId": "1743fa0000001000000000000000002f0000000000000000","ParentId": "1743fbb70000100000000000000000310000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Events": [{"level": "Debug","message": "debug string message"},{"level": "Debug","message": "debug field message"}],"metrics": [{"root thread": 0,"Attributes": {"root": "root span"},"Labels": ["root","metric"]}],"externalSpans": [{"TraceId": "","ParentId": "","InternalParentId": "17410dc600001000000000000000002d0000000000000000","SpanId": "1742fdc800001000000000000000002e0000000000000000","StartTime": -62135596800,"EndTime": -62135596800,"Attributes": {}},{"TraceId": "1743fb330000100000000000000000300000000000000000","ParentId": "1743fbb70000100000000000000000310000000000000000","InternalParentId": "1743fa0000001000000000000000002f0000000000000000","SpanId": "17440a420000100000000000000000330000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Attributes": {"method": "test","host": "test","attr0": "1743fb330000100000000000000000300000000000000000","attr1": "1743fbb70000100000000000000000310000000000000000","attr2": "1743fbb70000100000000000000000310000000000000000","attr3": "17440a420000100000000000000000330000000000000000"}}]}

{"TraceId": "1743fb330000100000000000000000300000000000000000","SpanId": "174402580000100000000000000000320000000000000000","ParentId": "1743fa0000001000000000000000002f0000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Events": [{"level": "Debug","message": "debug string"}],"metrics": [],"externalSpans": []}
```

## 真实数据样例

在使用第三方库时，由于第三方库日志存在输出到标准输出的情况，因此真实日志数据中会存在脏数据，真实数据如下，处理方式查看第一节数据分层。

```text
{"TraceId": "1743fb330000100000000000000000300000000000000000","SpanId": "1743fa0000001000000000000000002f0000000000000000","ParentId": "1743fbb70000100000000000000000310000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Events": [{"level": "Debug","message": "debug string message"},{"level": "Debug","message": "debug field message"}],"metrics": [{"root thread": 0,"Attributes": {"root": "root span"},"Labels": ["root","metric"]}],"externalSpans": [{"TraceId": "","ParentId": "","InternalParentId": "17410dc600001000000000000000002d0000000000000000","SpanId": "1742fdc800001000000000000000002e0000000000000000","StartTime": -62135596800,"EndTime": -62135596800,"Attributes": {}},{"TraceId": "1743fb330000100000000000000000300000000000000000","ParentId": "1743fbb70000100000000000000000310000000000000000","InternalParentId": "1743fa0000001000000000000000002f0000000000000000","SpanId": "17440a420000100000000000000000330000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Attributes": {"method": "test","host": "test","attr0": "1743fb330000100000000000000000300000000000000000","attr1": "1743fbb70000100000000000000000310000000000000000","attr2": "1743fbb70000100000000000000000310000000000000000","attr3": "17440a420000100000000000000000330000000000000000"}}]}

I Am Faker. ahhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh.

{"TraceId": "1743fb330000100000000000000000300000000000000000","SpanId": "174402580000100000000000000000320000000000000000","ParentId": "1743fa0000001000000000000000002f0000000000000000","StartTime": 1620807308,"EndTime": 1620807308,"Events": [{"level": "Debug","message": "debug string"}],"metrics": [],"externalSpans": []}