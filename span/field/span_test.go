package field

import (
	"sync"
	"testing"
	"time"

	"gotest.tools/assert"
)

func testFaterSonRelation(t *testing.T, f, c InternalSpan) {
	assert.Equal(t, f.ID(), c.ParentID())
	assert.Equal(t, f.TraceID(), c.TraceID())
}

func TestGenSpanID(t *testing.T) {
	total := 1000000
	store := map[string]int{}
	for i := 0; i < total; i++ {
		id := GenSpanID()
		if _, ok := store[id]; ok {
			t.Error("SpanID repeat")
			t.FailNow()
		} else {
			store[id] = 0
		}

	}
}

func TestGenTraceID(t *testing.T) {
	total := 1000000
	store := map[string]int{}
	for i := 0; i < total; i++ {
		id := GenTraceID()
		if _, ok := store[id]; ok {
			t.Error("SpanID repeat")
			t.FailNow()
		} else {
			store[id] = 0
		}

	}
}

func TestNewSpanFromPool(t *testing.T) {
	s0 := NewSpanFromPool(func(InternalSpan) {}, "")
	s1 := NewSpanFromPool(func(InternalSpan) {}, "")

	assert.Assert(t, s0 != s1, "Should get the same span")

	s0.Free()
	s01 := NewSpanFromPool(func(InternalSpan) {}, "")
	assert.Assert(t, s0 == s01, "should get same span")
}

func TestChildren(t *testing.T) {
	traceID := "test trace"
	s0 := NewSpanFromPool(func(InternalSpan) {}, traceID)
	assert.Equal(t, traceID, s0.TraceID())
	sc0 := s0.Children()
	sc1 := s0.Children()
	testFaterSonRelation(t, s0, sc0)
	testFaterSonRelation(t, s0, sc1)
	assert.Equal(t, 2, len(s0.ListChildren()))
}

func TestSignal(t *testing.T) {
	lock := &sync.Mutex{}
	count := 0
	s0 := NewSpanFromPool(func(s InternalSpan) {
		lock.Lock()
		count += 1
		lock.Unlock()
		s.Free()
	}, "")

	s0.Signal()

	s1 := NewSpanFromPool(func(s InternalSpan) {
		lock.Lock()
		count += 1
		lock.Unlock()
		s.Free()
	}, "")

	s1.Signal()

	start := time.Now()

	s2 := NewSpanFromPool(func(InternalSpan) {
		cost := time.Since(start)
		assert.Assert(t, cost < 1*time.Microsecond, "parent span signal() complete before children span")
	}, "")

	sc20 := s2.Children()
	s2.Signal()

	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 2, count)

	sc20.Signal()
	s2.Free()

}

func TestSetTraceID(t *testing.T) {
	s0 := NewSpanFromPool(func(InternalSpan) {}, "test trace")
	id := GenTraceID()
	s0.SetTraceID(id)
	assert.Equal(t, id, s0.TraceID())
}

func TestNewExternalSpan(t *testing.T) {
	exID := GenSpanID()
	traceID := GenTraceID()
	s0 := NewSpanFromPool(func(InternalSpan) {}, traceID)
	s0.SetParentID(exID)
	es0 := s0.NewExternalSpan()

	assert.Equal(t, es0.traceID, traceID)
	assert.Equal(t, es0.parentID, exID)
	assert.Equal(t, es0.internalParentID, s0.ID())

	sc0 := s0.Children()
	es1 := sc0.NewExternalSpan()
	assert.Equal(t, es1.traceID, traceID)
	assert.Equal(t, es1.parentID, exID)
	assert.Equal(t, es1.internalParentID, sc0.ID())

	sc0.NewExternalSpan()
	assert.Equal(t, 2, len(sc0.ListExternalSpan()))

}

func TestInternalSpanRecord(t *testing.T) {
	traceID := GenTraceID()
	s0 := NewSpanFromPool(func(InternalSpan) {}, traceID)

	now := time.Now()
	arrayField := MallocArrayField(4)
	arrayField.Append(IntField(1))
	arrayField.Append(Float64Field(1))
	arrayField.Append(StringField("test string in array"))
	arrayField.Append(TimeField(now))

	structField := MallocStructField(4)
	structField.Set("int", IntField(2))
	structField.Set("float", Float64Field(2))
	structField.Set("string", StringField("test string in struct"))
	structField.Set("time", TimeField(now))

	records := []Field{
		IntField(0),
		Float64Field(0),
		StringField("test string"),
		TimeField(now),
		arrayField,
		structField,
	}

	for _, f := range records {
		s0.Record(f)
	}

	dst := s0.ListRecord()
	for i, v := range dst {
		assert.Equal(t, records[i], v)
	}
}

func TestInternalSpanMetric(t *testing.T) {
	traceID := GenTraceID()
	s0 := NewSpanFromPool(func(InternalSpan) {}, traceID)

	m := Mmetric{}
	m.Set("test metric", 0.0)
	m1 := Mmetric{}
	m1.Set("test metric 1", 1.0)
	s0.Metric(m)
	s0.Metric(m1)

	metrics := s0.ListMetric()
	assert.Equal(t, len(metrics), 2)
	assert.Equal(t, metrics[0].(*Mmetric).Name, m.Name)
	assert.Equal(t, metrics[1].(*Mmetric).Name, m1.Name)
	assert.Equal(t, metrics[0].(*Mmetric).Value, m.Value)
	assert.Equal(t, metrics[1].(*Mmetric).Value, m1.Value)

}
