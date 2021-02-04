package deadlock

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/petermattis/goid"
)

type rwlockCtx struct {
	gid      int64
	locktime int64
	ref      int64
}

type RWMutex struct {
	L           sync.RWMutex
	waitReaders map[int64]*rwlockCtx // 读锁抢不到的记录
	waitWriters map[int64]*rwlockCtx
	readCount   int64
	writer      lockCtx
	mu          sync.Mutex
}

func (m *RWMutex) beforeLock(gid int64) bool {
	m.mu.Lock()
	// 检查写锁是否已被抢占
	if m.writer.gid == gid {
		m.mu.Unlock()
		Opts.OnDeadlock()
		return true
	}
	if m.waitWriters == nil {
		m.waitWriters = make(map[int64]*rwlockCtx)
	}
	// 检查写锁是否重复等待
	wwrs := m.waitWriters[gid]
	if wwrs != nil {
		m.mu.Unlock()
		Opts.OnDeadlock()
		return true
	}
	// 记录写锁抢占
	m.waitWriters[gid] = &rwlockCtx{
		gid:      gid,
		locktime: time.Now().Unix(),
		ref:      1,
	}
	m.mu.Unlock()
	return false
}

func (m *RWMutex) Lock() {
	if Opts.Disable {
		m.L.Lock()
		return
	}

	gid := goid.Get()
	if m.beforeLock(gid) {
		return
	}
	// 加写锁
	m.L.Lock()
	m.afterLock(gid)
}

func (m *RWMutex) afterLock(gid int64) {
	// 移除写锁等待记录
	m.mu.Lock()
	if m.waitWriters == nil {
		m.waitWriters = make(map[int64]*rwlockCtx)
	}
	wwrs := m.waitWriters[gid]
	if wwrs != nil {
		wwrs.ref--
		if wwrs.ref <= 0 {
			delete(m.waitWriters, gid)
		}
	}
	// 记录独占写锁信息
	m.writer.gid = gid
	m.writer.locktime = time.Now().Unix()
	if m.readCount == 0 {
		detector.addRWLocker(m)
	}
	m.mu.Unlock()
}

func (m *RWMutex) Unlock() {
	if Opts.Disable {
		m.L.Unlock()
		return
	}
	// 解锁前处理
	m.mu.Lock()
	m.writer.gid = 0
	m.writer.locktime = 0
	if m.readCount == 0 {
		detector.delRWLocker(m)
	}
	m.mu.Unlock()
	// 解写锁
	m.L.Unlock()
}

func (m *RWMutex) beforeRLock(gid int64) {
	// 记录等待读锁抢占
	m.mu.Lock()
	if m.waitReaders == nil {
		m.waitReaders = make(map[int64]*rwlockCtx)
	}
	if m.waitReaders[gid] == nil {
		m.waitReaders[gid] = &rwlockCtx{
			gid:      gid,
			ref:      0,
		}
	}
	m.waitReaders[gid].ref++
	m.waitReaders[gid].locktime = time.Now().Unix()
	m.mu.Unlock()
}

func (m *RWMutex) RLock() {
	if Opts.Disable {
		m.L.RLock()
		return
	}
	gid := goid.Get()
	m.beforeRLock(gid)
	// 加读锁
	m.L.RLock()
	m.afterRLock(gid)
}

func (m *RWMutex) afterRLock(gid int64) {
	m.mu.Lock()
	if m.waitReaders == nil {
		m.waitReaders = make(map[int64]*rwlockCtx)
	}
	wrrs := m.waitReaders[gid]
	if wrrs != nil {
		wrrs.ref--
		if wrrs.ref <= 0 {
			delete(m.waitReaders, gid)
		}
	}
	m.readCount++
	if m.readCount == 1 && m.writer.gid == 0 {
		detector.addRWLocker(m)
	}
	m.mu.Unlock()
}

func (m *RWMutex) RUnlock() {
	if Opts.Disable {
		m.L.RUnlock()
		return
	}
	// 解锁前处理
	m.mu.Lock()
	m.readCount--
	if m.readCount == 0 && m.writer.gid == 0 {
		detector.delRWLocker(m)
	}
	m.mu.Unlock()
	// 解读锁
	m.L.RUnlock()
}

func (m *RWMutex) RLocker() sync.Locker {
	return m
}

func (m *RWMutex) detect(now, timeout int64, stacks map[int64][]byte) {
	m.mu.Lock()
	deadGids := make(map[int64]bool)
	if m.writer.gid != 0 && m.writer.locktime+timeout <= now {
		deadGids[m.writer.gid] = true
	}
	for gid, ctx := range m.waitReaders {
		if ctx.locktime+timeout <= now {
			deadGids[gid] = true
		}
	}
	for gid, ctx := range m.waitWriters {
		if ctx.locktime+timeout <= now {
			deadGids[gid] = true
		}
	}
	m.mu.Unlock()

	if len(deadGids) == 0 {
		return
	}

	gNum := uint32(0)
	Opts.Logger.Write(logHeader)
	fmt.Fprintf(Opts.Logger, "the lock %p was grabbed by goroutines:\n", m)
	for gid := range deadGids {
		gNum++
		fmt.Fprintf(Opts.Logger, "[%d]\n", gid)
		if stacks[gid] != nil && isValidPrintStack(gNum) {
			Opts.Logger.Write(stacks[gid])
		}
	}
	Opts.Logger.Write([]byte("\nother waitting goroutines:\n"))

	m.mu.Lock()
	for gid := range m.waitReaders {
		gNum++
		if stacks[gid] != nil && deadGids[gid] == false && isValidPrintStack(gNum) {
			Opts.Logger.Write(stacks[gid])
		}
	}
	for gid := range m.waitWriters {
		gNum++
		if stacks[gid] != nil && deadGids[gid] == false && isValidPrintStack(gNum) {
			Opts.Logger.Write(stacks[gid])
		}
	}
	m.mu.Unlock()

	if buf, ok := Opts.Logger.(*bufio.Writer); ok {
		buf.Flush()
	}
	Opts.OnDeadlock()
}
