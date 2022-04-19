package field

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

type LogSpan interface {
	GetAttributes() Field
	// Recode record log
	SetRecord(Field)

	GetRecord() Field

	TraceID() string
	SpanID() string

	GetLogLevel() Field
	SetLogLevel(field Field)

	// Signal notify parent span， this span's work is done
	// Span should do Signal after work end
	Signal()

	GetContext() context.Context

	SetOption(...LogOptionFunc)

	Free()
}

type attribute struct {
	Type    string
	Message Field
}

func NewAttribute(typ string, message Field) *attribute {
	return &attribute{
		Type:    typ,
		Message: message,
	}
}

type logSpanV1 struct {
	log      Field
	transfer func(LogSpan) // use to transfer span's ownership
	level    Field
	wg       sync.WaitGroup

	lock sync.RWMutex
	//traceID    string
	ctx context.Context

	genID      *randomIDGenerator
	attributes Field
}

var Pool = sync.Pool{
	New: func() interface{} {
		return newSpan(nil, context.Background())
	},
}

// NewSpan get span from sync.pool
func NewSpanFromPool(own func(LogSpan), ctx context.Context) LogSpan {
	s := Pool.Get().(*logSpanV1)
	// s.reset()
	s.transfer = own
	s.ctx = ctx
	return s
}

func newSpan(own func(LogSpan), ctx context.Context) LogSpan {
	s := &logSpanV1{}
	s.ctx = ctx
	s.init()
	s.transfer = own
	return s
}

func (l *logSpanV1) init() {
	l.wg = sync.WaitGroup{}
	l.lock = sync.RWMutex{}
	l.genID = defaultIDGenerator()
	l.level = StringField("Trace")
	l.reset()
}

func (l *logSpanV1) GetContext() context.Context {
	return l.ctx
}

func (l *logSpanV1) getTraceSpan() trace.Span {
	return trace.SpanFromContext(l.ctx)
}

func (l *logSpanV1) reset() {
	l.log = nil
	l.attributes = nil
	l.transfer = nil
	l.ctx = nil
}

func (l *logSpanV1) GetAttributes() Field {
	if l.attributes == nil {
		return MallocStructField(0)
	}
	return l.attributes
}

func (l *logSpanV1) SetOption(options ...LogOptionFunc) {
	if options == nil {
		return
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(l)
	}
}

func (l *logSpanV1) Signal() {
	if l.transfer != nil {
		l.transfer(l)
	}
}

func (l *logSpanV1) Free() {
	l.reset()
	Pool.Put(l)
}

func (l *logSpanV1) SetRecord(r Field) {
	l.log = r
}

func (l *logSpanV1) GetRecord() Field {
	return l.log
}

func (l *logSpanV1) GetLogLevel() Field {
	return l.level
}

func (l *logSpanV1) SetLogLevel(level Field) {
	l.level = level
}

// ctx is nil or not trace context,return true
func (l *logSpanV1) IsNilContext() bool {
	spanCtx := l.getTraceSpan().SpanContext()
	return l.ctx == nil || (spanCtx.TraceID() == trace.TraceID{} && spanCtx.SpanID() == trace.SpanID{})
}

func (l *logSpanV1) TraceID() string {
	if l.IsNilContext() {
		return l.genID.NewTraceID()
	}
	return l.getTraceSpan().SpanContext().TraceID().String()
}

func (l *logSpanV1) SpanID() string {
	if l.IsNilContext() {
		return l.genID.NewSpanID()
	}
	return l.getTraceSpan().SpanContext().SpanID().String()
}

type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

func (gen *randomIDGenerator) NewSpanID() string {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	gen.randSource.Read(sid[:])
	return sid.String()
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewTraceID() string {
	gen.Lock()
	defer gen.Unlock()
	tid := trace.TraceID{}
	gen.randSource.Read(tid[:])
	return tid.String()
}

func defaultIDGenerator() *randomIDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}
