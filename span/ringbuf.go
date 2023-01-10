package span

import (
	"sync"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"

	"golang.org/x/sync/semaphore"
)

type ringBuffer struct { //nolint
	cap        int
	capLocker  *semaphore.Weighted
	length     int
	head       int
	tail       int
	buffer     []field.LogSpan
	headLocker *sync.RWMutex
	tailLocker *sync.RWMutex
	dataLocker sync.Locker
}

func initRingBuffer(l int) ringBuffer { //nolint

	return ringBuffer{
		cap:        l,
		length:     0,
		head:       0,
		tail:       0,
		buffer:     make([]field.LogSpan, l),
		headLocker: &sync.RWMutex{},
		tailLocker: &sync.RWMutex{},
		dataLocker: &sync.Mutex{},
		capLocker:  semaphore.NewWeighted(int64(l)),
	}
}

func (b *ringBuffer) push(s field.LogSpan) { //nolint
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

func (b *ringBuffer) pull() field.LogSpan { //nolint
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
