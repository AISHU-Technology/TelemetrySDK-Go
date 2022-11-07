# 接口文档

// Tracer 是一个全局变量，用于在业务代码中生产Span。

```
artrace.Tracer.Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span)
```

// SetAttributes 用于在Span中添加和业务相关的有语义的一系列键值对。

```
func SetAttributes(kv ...attribute.KeyValue)
```

// AddEvent 用于在Span中添加某一时刻发生的有意义的事件。

```
func AddEvent(name string, options ...EventOption)
```

// GetResource 获取内置资源信息，记录客户服务名，需要传入服务名 serviceName ，服务版本 serviceVersion ，服务实例ID。

```
func GetResource(serviceName string, serviceVersion string, serviceInstanceID string) *resource.Resource
```

// NewExporter 新建Exporter，需要传入指定的数据发送客户端 client.Client 。

```
func NewExporter(c client.Client) *client.Exporter
```

// NewHTTPClient 创建 client.Exporter 需要的HTTP数据发送客户端。

```
func NewHTTPClient(opts ...config.HTTPOption) client.Client
```

// NewStdoutClient 创建 client.Exporter 需要的Local数据发送客户端。

```
func NewStdoutClient(stdoutPath string) client.Client
```

// WithAnyRobotURL 设置 client.httpClient 数据上报地址。

```
func WithAnyRobotURL(URL string) config.Option
```

// WithCompression 设置Trace压缩方式：0代表无压缩，1代表GZIP压缩。

```
func WithCompression(compression int) config.Option
```

// WithHeader 设置 client.httpClient 用户自定义请求头。

```
func WithHeader(headers map[string]string) config.Option
```

// WithRetry 设置 client.httpClient 重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。

```
func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, ...) config.Option
```

// WithTimeout 设置 client.httpClient 连接超时时间。

```
func WithTimeout(duration time.Duration) config.Option
```

// ForceFlush 立即发送之前生产的链路数据。

```
func (p *TracerProvider) ForceFlush(ctx context.Context) error
```

// RegisterSpanProcessor 注册 SpanProcessor 来生产发送链路数据。

```
func (p *TracerProvider) RegisterSpanProcessor(s SpanProcessor)
```

// UnregisterSpanProcessor 解绑 SpanProcessor 来停止生产发送链路数据。

```
func (p *TracerProvider) UnregisterSpanProcessor(s SpanProcessor)
```