package field

import (
    // "span/field"
    "encoding/binary"
    "encoding/hex"
    "sync"
    "time"
)

var BaseNum uint64
var BaseNumLock sync.Locker
var ServiceID [8]byte

func init() {
    BaseNum = uint64(0)
    ServiceID = [8]byte{}
    BaseNumLock = &sync.RWMutex{}
}

// type OpenTelemetry struct {
//  Time       time.Time
//  Trace      []*HttpInfo
//  Severity   int
//  Name       string
//  Body       interface{}
//  Resource   interface{}
//  Attributes interface{}
//  Sub        []OpenTelemetry
// }

// GenSpanID is unsafe
func GenSpanID() string {
    BaseNumLock.Lock()
    defer BaseNumLock.Unlock()
    BaseNum += 1
    return genSpanID(ServiceID, BaseNum)

}

func GenTraceID() string {
    return GenSpanID()
}

// genSpanID support 2^64 item/ns, thread unsafe
func genSpanID(service [8]byte, base uint64) string {
    var id [24]byte
    now := time.Now().Nanosecond()
    timeLow := uint32(now & 0xffffffff)
    timeMid := uint16((now >> 32) & 0xffff)
    timeHi := uint16((now >> 48) & 0x0fff)
    timeHi |= 0x1000 // Version 1

    binary.BigEndian.PutUint32(id[0:], timeLow)
    binary.BigEndian.PutUint16(id[4:], timeMid)
    binary.BigEndian.PutUint16(id[6:], timeHi)

    base += 1
    binary.BigEndian.PutUint64(id[8:], base)
    copy(id[16:], service[:])

    return hex.EncodeToString(id[:])

}

type InternalSpan interface {
    // Children create children span, children work for sub thread or sub task
    Children() InternalSpan
    ListChildren() []InternalSpan
    // Parent(Span)
    ParentID() string
    SetParentID(string)

    // SetAttributes
    SetAttributes(string, Field)
    GetAttributes() Field

    // Recode record log
    Record(Field)
    ListRecord() []Field

    // metric record metric info
    Metric(Mmetric)
    ListMetric() []Field

    // trace init a trace for request
    TraceID() string
    ID() string
    SetTraceID(ID string)
    ListExternalSpan() []Field

    // NewExternalSpan create a httpspan
    // user should write data to span before span.Signal() for thread safe
    NewExternalSpan() *ExternalSpanField

    // Signal notify parent span， this span's work is done
    // Span should do Signal after work end
    Signal()

    Time() time.Time
    Free()
}

type internalSpanV1 struct {
    Records             []Field
    Metrics             []Field
    httpSpan            []Field
    transfer            func(InternalSpan) // use to transfer span's ownership
    children            []InternalSpan
    id                  string
    wg                  sync.WaitGroup
    time                time.Time
    lock                sync.RWMutex
    traceID             string
    parentID            string
    externalParentField string
    attributes          Field
}

var Pool = sync.Pool{
    New: func() interface{} {
        return newSpan(nil, "")
    },
}

// NewSpan get span from sync.pool
func NewSpanFromPool(own func(InternalSpan), traceID string) InternalSpan {
    s := Pool.Get().(*internalSpanV1)
    // s.reset()
    s.transfer = own
    s.time = time.Now()
    s.id = GenSpanID()
    s.traceID = traceID
    return s
}

func newSpan(own func(InternalSpan), traceID string) InternalSpan {
    s := &internalSpanV1{}
    s.init()
    s.id = GenSpanID()
    s.transfer = own
    s.traceID = traceID
    return s
}

func (l *internalSpanV1) init() {
    l.time = time.Now()

    l.id = GenSpanID()
    l.wg = sync.WaitGroup{}
    l.lock = sync.RWMutex{}
    l.reset()
}

func (l *internalSpanV1) reset() {
    l.httpSpan = l.httpSpan[:0]
    l.Records = l.Records[:0]
    l.Metrics = l.Metrics[:0]
    l.children = l.children[:0]
    l.transfer = nil
}

func (l *internalSpanV1) SetAttributes(t string, attrs Field) {
    a := MallocStructField(2)
    a.Set("type", StringField(t))
    a.Set("Attributes", attrs)
    l.attributes = a
}

func (l *internalSpanV1) GetAttributes() Field {
    return l.attributes
}

func (l *internalSpanV1) Children() InternalSpan {
    return l.newChildren()
}

func (l *internalSpanV1) newChildren() *internalSpanV1 {
    l.wg.Add(1)
    s := NewSpanFromPool(func(InternalSpan) {
        l.wg.Done()
    }, l.traceID).(*internalSpanV1)
    s.SetParentID(l.id)
    s.externalParentField = l.externalParentField
    l.children = append(l.children, s)
    s.attributes = l.attributes
    return s
}

func (l *internalSpanV1) ListChildren() []InternalSpan {
    return l.children[:]
}

func (l *internalSpanV1) Signal() {
    go func() {
        l.wg.Wait()
        if l.transfer != nil {
            l.transfer(l)
        }
    }()
    // l.wg.Wait()
    // if l.transfer != nil {
    //  l.transfer(l)
    // }
}

func (l *internalSpanV1) Free() {
    for _, c := range l.children {
        c.Free()
    }
    l.reset()
    Pool.Put(l)
}

func (l *internalSpanV1) Time() time.Time {
    return l.time
}

func (l *internalSpanV1) ID() string {
    return l.id
}

func (l *internalSpanV1) Record(r Field) {
    l.Records = append(l.Records, r)
}

func (l *internalSpanV1) ListRecord() []Field {
    return l.Records[:]
}

func (l *internalSpanV1) Metric(m Mmetric) {
    l.Metrics = append(l.Metrics, &m)
}

func (l *internalSpanV1) ListMetric() []Field {
    return l.Metrics[:]
}

func (l *internalSpanV1) TraceID() string {
    return l.traceID

}

func (l *internalSpanV1) SetTraceID(ID string) {
    l.traceID = ID
}

func (l *internalSpanV1) ParentID() string {
    return l.parentID
}

func (l *internalSpanV1) SetParentID(ID string) {
    l.parentID = ID
    l.externalParentField = ID
}

// NewHttpSpan create a httpspan
// user should write data to span before span.signal
func (l *internalSpanV1) NewExternalSpan() *ExternalSpanField {
    span := &ExternalSpanField{
        traceID:          l.traceID,
        id:               GenSpanID(),
        parentID:         l.externalParentField,
        internalParentID: l.id,
    }
    l.httpSpan = append(l.httpSpan, span)
    return span
}

func (l *internalSpanV1) ListExternalSpan() []Field {
    return l.httpSpan[:]
}

