package mutex

import (
	"testing"
	"time"
	"reflect"
)

// Custom data for testing
type CriticalDataSection struct {
	field1 int
}

// Test mutex for successful getting
func TestMutexGet(t *testing.T) {
	mutex := Mutex{ state: UNLOCKED, timeout: DEFAULT_TIMEOUT }
	if !mutex.getMutex() {
		t.Error("Mutex can't be get")
	}
	mutex.releaseMutex()
	if !mutex.getMutex() {
		t.Error("Mutex can't be get after release")
	}
}

// Test mutex for successful releasing
func TestMutexRealese(t *testing.T) {
	mutex := Mutex{ state: UNLOCKED, timeout: DEFAULT_TIMEOUT }
	mutex.getMutex()
	if !mutex.releaseMutex() {
		t.Error("Mutex can't be released")
	}
	defer func() {
		if err := recover(); err == nil {
			t.Error("Mutex can't be released twice")
		}
	}()
	mutex.releaseMutex()
}

// Test mutex for throwing timeout
func TestMutexTimeOut(t *testing.T) {
	mutex := Mutex{ state: UNLOCKED, timeout: DEFAULT_TIMEOUT }
	mutex.getMutex()
	defer func() {
		if err := recover(); reflect.TypeOf(err).String() != "string" {
			t.Error("Mutex timeout doesn't work")
		}
	}()
	mutex.getMutex()
}

// Test mutex for solving race conditions
func TestMutexRaceCondition(t *testing.T) {
	data := &CriticalDataSection{field1: 0}
	mutex := Mutex{ state: UNLOCKED, timeout: DEFAULT_TIMEOUT }
	for i := 1; i < 10; i++ {
		go func(i int) {
			mutex.getMutex()
			data.field1 = i
			mutex.releaseMutex()
		}(i)
	}
	time.Sleep(1 * time.Second)
}
