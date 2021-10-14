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
	SetAttributes(attribute *Attribute)

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

	SetContext(ctx context.Context)

	GetContext() context.Context

	Free()
}

type Attribute struct {
	Type    string
	Message Field
}

func NewAttribute(typ string, message Field) *Attribute {
	return &Attribute{
		Type:    typ,
		Message: message,
	}
}

type LogSpanV1 struct {
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
	s := Pool.Get().(*LogSpanV1)
	// s.reset()
	s.transfer = own
	s.ctx = ctx
	return s
}

func newSpan(own func(LogSpan), ctx context.Context) LogSpan {
	s := &LogSpanV1{}
	s.ctx = ctx
	s.init()
	s.transfer = own
	return s
}

func (l *LogSpanV1) init() {
	l.wg = sync.WaitGroup{}
	l.lock = sync.RWMutex{}
	l.genID = defaultIDGenerator()
	l.reset()
}

func (l *LogSpanV1) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *LogSpanV1) GetContext() context.Context {
	return l.ctx
}

func (l *LogSpanV1) getTraceSpan() trace.Span {
	return trace.SpanFromContext(l.ctx)
}

func (l *LogSpanV1) reset() {
	l.log = nil
	l.attributes = nil
	l.transfer = nil
	l.ctx = nil
}

func (l *LogSpanV1) SetAttributes(attribute *Attribute) {
	record := MallocStructField(2)
	record.Set(attribute.Type, attribute.Message)
	record.Set("Type", StringField(attribute.Type))
	l.attributes = record
}

func (l *LogSpanV1) GetAttributes() Field {
	if l.attributes == nil {
		return MallocStructField(0)
	}
	return l.attributes
}

func (l *LogSpanV1) Signal() {
	go func() {
		l.wg.Wait()
		if l.transfer != nil {
			l.transfer(l)
		}
	}()

}

func (l *LogSpanV1) Free() {
	l.reset()
	Pool.Put(l)
}

func (l *LogSpanV1) SetRecord(r Field) {
	l.log = r
}

func (l *LogSpanV1) GetRecord() Field {
	return l.log
}

func (l *LogSpanV1) GetLogLevel() Field {
	return l.level
}

func (l *LogSpanV1) SetLogLevel(level Field) {
	l.level = level
}

// ctx is nil or not trace context,return true
func (l *LogSpanV1) IsNilContext() bool {
	spanCtx := l.getTraceSpan().SpanContext()
	return l.ctx == nil || (spanCtx.TraceID() == trace.TraceID{} && spanCtx.SpanID() == trace.SpanID{})
}

func (l *LogSpanV1) TraceID() string {
	if l.IsNilContext() {
		return l.genID.NewTraceID()
	}
	return l.getTraceSpan().SpanContext().TraceID().String()
}

func (l *LogSpanV1) SpanID() string {
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
