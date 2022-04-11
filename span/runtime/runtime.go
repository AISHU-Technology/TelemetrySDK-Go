package runtime

import (
	"context"
	"sync"

	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/field"
	"devops.aishu.cn/AISHUDevOps/AnyRobot/_git/DE_TelemetryGo.git/span/open_standard"
)

// Runtime read data from channel and write data in a single goroutine
type Runtime struct {
	cache   chan field.LogSpan
	builder func(func(field.LogSpan), context.Context) field.LogSpan
	stop    chan int
	wg      *sync.WaitGroup
	// id      uint64
	close     bool
	closeLock sync.RWMutex
	runLock   sync.Mutex
	w         open_standard.Writer
	once      sync.Once
}

// NewRuntime return a runtime
func NewRuntime(w open_standard.Writer, builder func(func(field.LogSpan), context.Context) field.LogSpan) *Runtime {
	r := &Runtime{
		cache:   make(chan field.LogSpan, 100),
		builder: builder,
		stop:    make(chan int, 1),
		wg:      &sync.WaitGroup{},
		// false = running, true = closed or closing
		close: false,
		// protect the close value
		closeLock: sync.RWMutex{},
		// represent the state of runtime thread, runtime thread will lock this until thread over
		runLock: sync.Mutex{},
		w:       w,
		once:    sync.Once{},
	}

	return r
}

// Children return a logger span
// if Runtime has been close return nil
// user should return span's onwership after Span is useless by Span.Signal()
func (r *Runtime) Children(ctx context.Context) field.LogSpan {
	// TODO: remove read lock
	r.closeLock.RLock()
	defer r.closeLock.RUnlock()

	if r.close {
		return nil
	}
	s := r.builder(r.transfer, ctx)
	r.wg.Add(1)
	return s
}

//Signal stop runtime thread
func (r *Runtime) Signal() {
	r.closeLock.Lock()
	r.close = true
	r.wg.Wait()
	r.once.Do(func() {
		r.stop <- 0
		close(r.cache)
		// r.close = true
		r.w.Close()
	})

	r.runLock.Lock()
	r.runLock.Unlock()
	r.closeLock.Unlock()
}

func (r *Runtime) transfer(s field.LogSpan) {
	r.cache <- s
}

// Run will deal Runtime's span in current go runtine
// Run will block the go runtine, so you may need to use a go runtine to run
// Runtime will close when all span has been did when receive system signal or do Signal()
// Maybe runtime should not close by signal
func (r *Runtime) Run() {
	r.runLock.Lock()
	defer r.runLock.Unlock()
	for {
		s, ok := <-(r.cache)
		if ok != true {
			err := r.w.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		r.w.Write(s)
		s.Free()
		r.wg.Done()
	}
}
