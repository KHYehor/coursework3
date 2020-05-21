package mutex

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

const (
	LOCKED = 1
	UNLOCKED = 0
	ERROR_GET = "Gorutine holds Mutex longer than: %d ms"
	ERROR_RELEASE = "You can not release mutex because it is not grabbed"
	DEFAULT_TIMEOUT = time.Second * 5
)

// Mutex with timeout
type Mutex struct {
	state int32
	timeout time.Duration
}

// Getting mutex for control
func (m *Mutex) getMutex() bool {
	if atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED) {
		return true
	}
	// Waiting for getting mutex
	start := time.Now()
	for {
		if atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED) {
			// Finish stopwatch
			total := time.Now().Sub(start)
			// Printing total time of waiting for the mutex
			log.Printf("Mutex has been holding: %d ms", total.Microseconds())
			return true
		}
		total := time.Now().Sub(start)
		if total > m.timeout {
			panic(fmt.Sprintf(ERROR_GET, m.timeout.Microseconds()))
		}
	}
}

// Releasing mutex to let deal with it
func (m *Mutex) releaseMutex() bool {
	if atomic.CompareAndSwapInt32(&m.state, LOCKED, UNLOCKED) {
		return true
	}
	panic(fmt.Sprintf(ERROR_RELEASE))
}
