package log

import (
	"os"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/encoder"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/runtime"
	"time"
)

/*
As a quick start. user can log info without InternalSpan. it's easy to use.
when log info without InternalSpan, user should pass span parameter as nil like follow operator.
if user don't use InternalSpan, their log infos can't form a trace.
The log's traceID and SpanID is meaningless.
And Some logger's interfaces should't be use. if user call the interface, it will do nothing.
*/
func Example_simpleDoc() {
	/*
		As a quick start. user can log info without InternalSpan.
		when log info without InternalSpan, user should pass span parameter as nil like follow operator.
		if user don't use InternalSpan, their log infos can't form a trace.
		The log's traceID and SpanID is meaningless.
		And Some logger's interfaces should't be use. if user call the interface, it will do nothing.

		First we need to create a gloabl logger by l := easylog.NewdefaultSamplerLogger().

		we can store string log or object log without InternalSpan like op 1.1. the info will output immediately.

		we can store some metrics to InternalSpan like op 1.2 .

		if we want to log a sub task or a sub thread info. we can do like 2.x .

		Please always remember the two most important points:

				* Rememeber to signal the InternalSpan after it's useless. The parent InternalSpan will wait child InternalSpan to be signal
				* InternalSpan's operator is thread unsafe. One thread, one InternalSpan
	*/

	// get Default logger
	l := NewdefaultSamplerLogger()

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

	// 1.0 log message into roor internalSpan
	l.Debug("debug string message", nil)
	l.DebugField(field.StringField("debug field message"), "test", nil)

	type A struct {
		Name string
		Age int
	}
	var a = &A{Name:"123",Age:12,}


	l.DebugField(field.MallocJsonFiled(a),"detail",nil)

	// 1.1 start a new thread for sub task
	go func() {
		// 2.1 log message  for child thread
		l.Debug("debug string", nil)

	}()

	// 1.2 record some metric
	m := field.Mmetric{}
	m.Set("root thread", 0.0)
	m.AddLabel("root")
	m.AddLabel("metric")
	m.AddAttribute("root", "root span")
	l.RecordMetrics(m, nil)

	// final close runtime and clean work space
	l.Close()

}

/*
As a quick start with InternalSpan. This case can use all interface of logger.
It is more difficult to use, but It can log more plump info and do a link trace .
*/
func Example_allDoc() {
	/*
		As a quick start

		First we need to create a gloabl logger by l := easylog.NewdefaultSamplerLogger().

		Then we need to create a root InteranlSpan to store log info.

		we can store string log or object log to InternalSpan like op 1.1 .

		we can set some attributes to InternalSpan like op 1.2. I suggest only set attributes in root InternalSpan.
		because we set attributes to child InternalSpan, the root can't dispaly child's attributes.

		we can store some metrics to InternalSpan like op 1.5 .

		if we need to do some call  out of this process, we can use ExternalSpan to record call info like op 1.6 .

		if we won't use a InternalSpan, we should signal it. It's very import thing. Logger will output the InternalSpan
		when the root InternalSpan is signal. And The InternalSpan will wait its chilid InternalSpan signal.

		if we wan't to log a sub task or a sub thread info. we should create a child InternalSpan.
		from InternalSpan before InternalSpan be signal like op 1.3 .
		after that we can use child Internal Span in sub thread/task like op 2.X .
		Finally, don't forget to signal InternalSpan after it is uselese.

		Please always remember the two most important points:

				* Rememeber to signal the InternalSpan after it's useless. The parent InternalSpan will wait child InternalSpan to be signal.
				* InternalSpan's operator is thread unsafe. One thread, one InternalSpan.
	*/

	// get Default logger
	l := NewdefaultSamplerLogger()

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

	// 1. first create a root internalSpan
	root := l.NewInternalSpan()

	// 1.0 set trace info for root internalSpan
	traceID := field.GenSpanID()
	externalParentID := field.GenSpanID()
	l.SetTraceID(traceID, root)
	l.SetParentID(externalParentID, root)

	// 1.1 log message into roor internalSpan
	l.Debug("debug string message", root)
	l.DebugField(field.StringField("debug field message"), "test", root)

	// 1.2 set attributes for a span
	attrs := field.MallocStructField(3)
	attrs.Set("work", field.StringField("test"))
	attrs.Set("testFunc", field.StringField("TestSamplerLogger"))
	attrs.Set("testSpan", field.StringField("root"))
	l.SetAttributes("SampleLogerTest", attrs, root)

	// 1.3 create a child internalSpan from root for a sub thread/task
	child0 := l.ChildrenInternalSpan(root)

	// 1.4 start a new thread for sub task
	go func() {
		// 2.1 log message into child internalSpan for child thread
		l.Debug("debug string", child0)

		// 2.X signal child0
		child0.Signal()
	}()

	// 1.5 record some metric into root internalSpan
	m := field.Mmetric{}
	m.Set("root thread", 0.0)
	m.AddLabel("root")
	m.AddLabel("metric")
	m.AddAttribute("root", "root span")
	l.RecordMetrics(m, root)

	// 1.6 record first external request into root internalSpan
	es, err := l.NewExternalSpan(root)
	if err != nil {
		panic(err)
	}

	// 1.6.1 get trace info for some work
	tID := es.TraceID()
	espID := es.ParentID()
	parentID := es.ParentID()
	spanID := es.ID()
	// 1.6.2 write info to external span
	es.StartTime = time.Now()
	es.EndTime = time.Now()
	es.Attributes.Set("method", field.StringField("test"))
	es.Attributes.Set("host", field.StringField("test"))
	es.Attributes.Set("attr0", field.StringField(tID))
	es.Attributes.Set("attr1", field.StringField(espID))
	es.Attributes.Set("attr2", field.StringField(parentID))
	es.Attributes.Set("attr3", field.StringField(spanID))

	// 1.X signal root internalSpan
	root.Signal()

	// final close runtime and clean work space
	l.Close()
}
