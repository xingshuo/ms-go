package deadlock

import (
	"sync"
	"sync/atomic"
)
import "github.com/petermattis/goid"

type Mutex struct {
	L sync.Mutex
}

func (m *Mutex) Lock() {
	if Opts.Disable {
		m.L.Lock()
		return
	}
	gid := goid.Get()
	detector.addWaiter(m, gid)
	m.L.Lock()
	detector.addLocker(m, gid)
}

func (m *Mutex) Unlock() {
	if Opts.Disable {
		m.L.Unlock()
		return
	}
	m.L.Unlock()
	detector.delLocker(m)
}

type RWMutex struct {
	L sync.RWMutex
	Ref int64
	WaitRef int64
	lastGid int64
}

func (m *RWMutex) Lock() {
	if Opts.Disable {
		panic("lock disable")
		m.L.Lock()
		return
	}
	gid := goid.Get()
	detector.addWaiter(m, gid)
	atomic.AddInt64(&m.WaitRef, 1)
	m.L.Lock()
	m.lastGid = gid
	atomic.AddInt64(&m.Ref, 1)
	atomic.AddInt64(&m.WaitRef, -1)
	detector.addLocker(m, gid)
}

func (m *RWMutex) Unlock() {
	if Opts.Disable {
		panic("unlock disable")
		m.L.Unlock()
		return
	}
	atomic.AddInt64(&m.Ref, -1)
	m.lastGid = 0
	m.L.Unlock()
	detector.delLocker(m)
}

func (m *RWMutex) RLock() {
	if Opts.Disable {
		panic("rlock disable")
		m.L.RLock()
		return
	}
	gid := goid.Get()
	detector.addWaiter(m, gid)
	atomic.AddInt64(&m.WaitRef, 1)
	m.L.RLock()
	m.lastGid = gid
	atomic.AddInt64(&m.Ref, 1)
	atomic.AddInt64(&m.WaitRef, -1)
	detector.addLocker(m, gid)
}

func (m *RWMutex) RUnlock() {
	if Opts.Disable {
		panic("runlock disable")
		m.L.RUnlock()
		return
	}
	atomic.AddInt64(&m.Ref, -1)
	m.lastGid = 0
	m.L.RUnlock()
	detector.delLocker(m)
}

func (m *RWMutex) RLocker() sync.Locker {
	return m
}
