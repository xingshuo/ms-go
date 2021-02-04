package deadlock

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/petermattis/goid"
)

type Mutex struct {
	L       sync.Mutex
	waiters map[int64]int64 // gid : time.Now().Unix()
	owner   lockCtx         // gid
	mu      sync.Mutex
}

func (m *Mutex) Lock() {
	if Opts.Disable {
		m.L.Lock()
		return
	}
	gid := goid.Get()
	// 加锁前
	m.mu.Lock()
	if gid == m.owner.gid {
		m.mu.Unlock()
		Opts.OnDeadlock()
		return
	}
	if m.waiters == nil {
		m.waiters = make(map[int64]int64)
	}
	// 重复等锁
	if m.waiters[gid] > 0 {
		m.mu.Unlock()
		Opts.OnDeadlock()
		return
	}
	// 添加gid到waiter列表
	m.waiters[gid] = time.Now().Unix()
	m.mu.Unlock()
	// 加锁处理
	m.L.Lock()
	// 加锁后
	m.mu.Lock()
	delete(m.waiters, gid)
	m.owner.gid = gid
	m.owner.locktime = time.Now().Unix()
	detector.addLocker(m)
	m.mu.Unlock()
}

func (m *Mutex) Unlock() {
	if Opts.Disable {
		m.L.Unlock()
		return
	}
	// 解锁前
	m.mu.Lock()
	m.owner.gid = 0
	m.owner.locktime = 0
	detector.delLocker(m)
	m.mu.Unlock()
	// 解锁
	m.L.Unlock()
}

func (m *Mutex) detect(now, timeout int64, stacks map[int64][]byte) {
	m.mu.Lock()
	deadGids := make(map[int64]bool)
	if m.owner.gid != 0 && m.owner.locktime+timeout <= now {
		deadGids[m.owner.gid] = true
	}
	for gid, locktime := range m.waiters {
		if locktime+timeout <= now {
			deadGids[gid] = true
		}
	}
	m.mu.Unlock()

	if len(deadGids) == 0 {
		return
	}

	Opts.Logger.Write(logHeader)
	fmt.Fprintf(Opts.Logger, "the lock %p was grabbed by goroutines:\n", m)
	for gid := range deadGids {
		fmt.Fprintf(Opts.Logger, "[%d]\n", gid)
		if stacks[gid] != nil {
			Opts.Logger.Write(stacks[gid])
		}
	}
	Opts.Logger.Write([]byte("\nother waitting goroutines:\n"))

	m.mu.Lock()
	for gid := range m.waiters {
		if stacks[gid] != nil && deadGids[gid] == false {
			Opts.Logger.Write(stacks[gid])
		}
	}
	m.mu.Unlock()

	if buf, ok := Opts.Logger.(*bufio.Writer); ok {
		buf.Flush()
	}
	Opts.OnDeadlock()
}
