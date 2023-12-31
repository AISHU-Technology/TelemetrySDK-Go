package span

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"sync"

	"golang.org/x/sync/semaphore"
)

type ringBuffer struct {
	cap        int
	capLocker  *semaphore.Weighted
	length     int
	head       int
	tail       int
	buffer     []field.InternalSpan
	headLocker *sync.RWMutex
	tailLocker *sync.RWMutex
	dataLocker sync.Locker
}

func initRingBuffer(l int) ringBuffer {

	return ringBuffer{
		cap:        l,
		length:     0,
		head:       0,
		tail:       0,
		buffer:     make([]field.InternalSpan, l),
		headLocker: &sync.RWMutex{},
		tailLocker: &sync.RWMutex{},
		dataLocker: &sync.Mutex{},
		capLocker:  semaphore.NewWeighted(int64(l)),
	}
}

func (b *ringBuffer) push(s field.InternalSpan) {
	for {
		b.tailLocker.RLock()
		b.dataLocker.Lock()

		if b.length == b.cap {
			b.dataLocker.Unlock()
			b.tailLocker.RUnlock()
			continue
		}

		b.buffer[b.tail] = s
		b.tail = (b.tail + 1) % b.cap
		b.length += 1

		b.tailLocker.RUnlock()
		b.headLocker.Unlock()
		if b.length == b.cap {
			b.tailLocker.Lock()
		}
		b.dataLocker.Unlock()
		return

	}

}

func (b *ringBuffer) pull() field.InternalSpan {
	b.headLocker.RLock()
	defer b.headLocker.RUnlock()
	b.dataLocker.Lock()
	defer b.dataLocker.Unlock()

	r := b.buffer[b.length]
	b.length -= 1
	b.head = (b.head + 1) % b.cap

	b.tailLocker.Unlock()
	if b.length == 0 {
		b.headLocker.Lock()
	}
	return r
}
