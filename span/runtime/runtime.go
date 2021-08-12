package runtime

import (
	"sync"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/open_standard"
	// "go.uber.org/zap"
)

// var log *zap.Logger

func init() {
	// var err error
	// log, err = zap.NewProduction()

	// // err := log.NewConsoleWriter(false)
	// if err != nil {
	// 	panic(err)
	// }

}

// Runtime read data from channel and write data in a single goroutine
type Runtime struct {
	cache   chan field.InternalSpan
	builder func(func(field.InternalSpan), string) field.InternalSpan
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
func NewRuntime(w open_standard.Writer, builder func(func(field.InternalSpan), string) field.InternalSpan) *Runtime {
	r := &Runtime{
		cache:   make(chan field.InternalSpan, 100),
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
func (r *Runtime) Children() field.InternalSpan {
	// TODO: remove read lock
	r.closeLock.RLock()
	defer r.closeLock.RUnlock()

	if r.close {
		return nil
	}
	s := r.builder(r.transfer, "")
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
	})

	r.runLock.Lock()
	r.runLock.Unlock()
	r.closeLock.Unlock()
}

func (r *Runtime) transfer(s field.InternalSpan) {
	r.cache <- s
}

// Run will deal Runtime's span in current go runtine
// Run will block the go runtine, so you may need to use a go runtine to run
// Runtime will close when all span has been did when receive system signal or do Signal()
// Maybe runtime should not close by signal
func (r *Runtime) Run() {
	// start := time.Now().Second()
	// defer log.Sync()

	// sugar := log.Sugar()
	// writer := ioutil.Discard

	// use for generate sub span ID
	// support 2^64 span/ms

	// signalStop := make(chan os.Signal, 1)
	// signal.Notify(signalStop, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	// out := bufio.NewWriter(r.w)

	// b := bytes.NewBuffer(nil)

	// enc := encoder.NewJsonEncoder(r.w)

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
		// select {
		// // case <-signalStop:
		// // 	go r.Signal()
		// // case <-r.stop:
		// // 	close(r.cache)
		// case s, ok := <-(r.cache):
		// 	if !ok {
		// 		// sugar.Info("logger end, run: ", time.Now().Second()-start)
		// 		// fmt.Println("logger end, run: ", time.Now().Second()-start)
		// 		// r.w.Close()
		// 		// r.close = true
		// 		r.w.Close()
		// 		r.closeLock.Unlock()
		// 		return
		// 	}

		// 	// convertToOpenTelemetry1(enc, s)
		// 	r.w.Write(s)

		// 	// convertToOpenTelemetry(b, s)
		// 	// b.Reset()
		// 	// json.Marshal(s)
		// 	// b, _ := json.Marshal(s)
		// 	// writer.Write(b)
		// 	// record := convertToOpenTelemetry(s)
		// 	// data, err := json.Marshal(record)
		// 	// if err != nil && len(data) < 1 {
		// 	// 	fmt.Println(err)
		// 	// 	// sugar.Error(err)
		// 	// }
		// 	// log.Info(record)
		// 	// log.Info("", zap.Any("raw", record))
		// 	// sugar.Infow("test", record)
		// 	// fmt.Printf("%s\n", string(data))

		// 	s.Free()
		// 	r.wg.Done()
		// }
	}

}

