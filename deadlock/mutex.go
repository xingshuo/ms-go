package deadlock

import "sync"

type Mutex struct {
	L sync.Mutex
}

func (m *Mutex) Lock() {
	if !Opts.Disable {
		detector.addWaiter(m)
	}
	m.L.Lock()
	if !Opts.Disable {
		detector.addLocker(m)
	}
}

func (m *Mutex) Unlock() {
	m.L.Unlock()
	if !Opts.Disable {
		detector.delLocker(m)
	}
}

type RWMutex struct {
	L sync.RWMutex
}

func (m *RWMutex) Lock() {
	if !Opts.Disable {
		detector.addWaiter(m)
	}
	m.L.Lock()
	if !Opts.Disable {
		detector.addLocker(m)
	}
}

func (m *RWMutex) Unlock() {
	m.L.Unlock()
	if !Opts.Disable {
		detector.delLocker(m)
	}
}

func (m *RWMutex) RLock() {
	if !Opts.Disable {
		detector.addWaiter(m)
	}
	m.L.RLock()
	if !Opts.Disable {
		detector.addLocker(m)
	}
}

func (m *RWMutex) RUnlock() {
	m.L.RUnlock()
	if !Opts.Disable {
		detector.delLocker(m)
	}
}

func (m *RWMutex) RLocker() sync.Locker {
	return m
}
