package deadlock

import (
	"sync"
	"testing"
)

// mutex lock bench test
func BenchmarkRawMutex(b *testing.B) {
	var lock sync.Mutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkMutexDisable(b *testing.B) {
	Opts.Disable = true
	var lock Mutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkMutexEnable(b *testing.B) {
	Opts.Disable = false
	var lock Mutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

// rwmutex lock bench test
func BenchmarkRawRWMutex(b *testing.B) {
	var lock sync.RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkRWMutexDisable(b *testing.B) {
	Opts.Disable = true
	var lock RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkRWMutexEnable(b *testing.B) {
	Opts.Disable = false
	var lock RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}
