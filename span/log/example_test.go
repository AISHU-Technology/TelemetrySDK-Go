package log

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
	"os"
)

/*
As a quick start. user can log info without LogSpan. it's easy to use.
when log info without LogSpan, user should pass span parameter as nil like follow operator.
if user don't use LogSpan, their log infos can't form a trace.
The log's traceID and SpanID is meaningless.
And Some logger's interfaces should't be use. if user call the interface, it will do nothing.
*/
func Example_simpleDoc() {
	/*
		As a quick start. user can log info without LogSpan.
		when log info without LogSpan, user should pass span parameter as nil like follow operator.
		if user don't use LogSpan, their log infos can't form a trace.
		The log's traceID and SpanID is meaningless.
		And Some logger's interfaces should't be use. if user call the interface, it will do nothing.

		First we need to create a gloabl logger by l := easylog.NewdefaultSamplerLogger().

		we can store string log or object log without LogSpan like op 1.1. the info will output immediately.

		we can store some metrics to LogSpan like op 1.2 .

		if we want to log a sub task or a sub thread info. we can do like 2.x .

		Please always remember the two most important points:

				* Rememeber to signal the LogSpan after it's useless. The parent LogSpan will wait child LogSpan to be signal
				* LogSpan's operator is thread unsafe. One thread, one LogSpan
	*/

	// get Default logger
	l := NewDefaultSamplerLogger()

	// init Default logger
	output := os.Stdin
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
	}
	writer.SetDefultResources()
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	l.SetRuntime(run)

	// start runtime
	go run.Run()

	// Or we can use easylog.NewdefaultSamplerLogger to get SamplerLogger that can be used directly without init
	// like l := easylog.NewdefaultSamplerLogger

	// 1.0 log message into roor LogSpan
	l.Debug("debug string message", nil)
	l.DebugField(field.StringField("debug field message"), "test", nil)

	type A struct {
		Name string
		Age  int
	}
	var a = &A{Name: "123", Age: 12}

	l.DebugField(field.MallocJsonField(a), "detail", nil)

	// 1.1 start a new thread for sub task
	go func() {
		// 2.1 log message  for child thread
		l.Debug("debug string", nil)

	}()

	// final close runtime and clean work space
	l.Close()

}

/*
As a quick start with LogSpan. This case can use all interface of logger.
It is more difficult to use, but It can log more plump info and do a link trace .
*/
func Example_allDoc() {
	/*
		As a quick start

		First we need to create a gloabl logger by l := easylog.NewdefaultSamplerLogger().

		Then we need to create a root InteranlSpan to store log info.

		we can store string log or object log to LogSpan like op 1.1 .

		we can set some attributes to LogSpan like op 1.2. I suggest only set attributes in root LogSpan.
		because we set attributes to child LogSpan, the root can't dispaly child's attributes.

		we can store some metrics to LogSpan like op 1.5 .

		if we need to do some call  out of this process, we can use ExternalSpan to record call info like op 1.6 .

		if we won't use a LogSpan, we should signal it. It's very import thing. Logger will output the LogSpan
		when the root LogSpan is signal. And The LogSpan will wait its chilid LogSpan signal.

		if we wan't to log a sub task or a sub thread info. we should create a child LogSpan.
		from LogSpan before LogSpan be signal like op 1.3 .
		after that we can use child Internal Span in sub thread/task like op 2.X .
		Finally, don't forget to signal LogSpan after it is uselese.

		Please always remember the two most important points:

				* Rememeber to signal the LogSpan after it's useless. The parent LogSpan will wait child LogSpan to be signal.
				* LogSpan's operator is thread unsafe. One thread, one LogSpan.
	*/

	// get Default logger
	l := NewDefaultSamplerLogger()

	// init Default logger
	output := os.Stdin
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(output),
	}
	writer.SetDefultResources()
	run := runtime.NewRuntime(writer, field.NewSpanFromPool)
	l.SetRuntime(run)

	// start runtime
	go run.Run()

	// Or we can use easylog.NewdefaultSamplerLogger to get SamplerLogger that can be used directly without init
	// like l := easylog.NewdefaultSamplerLogger

	// 1.1 log message into roor LogSpan
	l.Debug("debug string message", nil)
	l.DebugField(field.StringField("debug field message"), "test", nil)

	// 1.2 set attributes for a span
	attrs := field.MallocStructField(3)
	attrs.Set("work", field.StringField("test"))
	attrs.Set("testFunc", field.StringField("TestSamplerLogger"))
	attrs.Set("testSpan", field.StringField("root"))

	a := field.NewAttribute("attr", attrs)
	l.Info("123", a)
	l.InfoField(field.StringField("112"), "123", a)

	// final close runtime and clean work space
	l.Close()
}
