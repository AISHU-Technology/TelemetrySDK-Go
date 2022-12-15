package runtime

import (
	"context"
	"sync"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/log_config"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/open_standard"
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
	Logs      []field.LogSpan
	Size      int
	Ticker    *time.Ticker
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
		Ticker:  time.NewTicker(log_config.Internal),
		Logs:    make([]field.LogSpan, 0, log_config.MaxLog+1),
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

// Signal stop runtime thread
func (r *Runtime) Signal() {
	r.closeLock.Lock()
	r.close = true
	r.wg.Wait()
	if r.Size > 0 {
		r.forceWrite()
	}
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
		select {
		case s, ok := <-(r.cache):
			// 关闭之后退出循环。
			if !ok {
				// 发完最后的数据关闭Exporter。
				err := r.w.Close()
				if err != nil {
					panic(err)
				}
				return
			}
			r.Logs = append(r.Logs, s)
			r.Size++
			r.wg.Done()
			// 超过上限发送。
			if r.Size >= log_config.MaxLog {
				r.forceWrite()
			}
		// 定时发送。
		case <-r.Ticker.C:
			if r.Size > 0 {
				r.forceWrite()
			}
		}
	}
}
func (r *Runtime) forceWrite() {
	r.w.Write(r.Logs)
	// 发送完之后清空队列。
	for _, v := range r.Logs {
		v.Free()
	}
	r.Size = 0
	r.Logs = make([]field.LogSpan, 0, log_config.MaxLog+1)
}
