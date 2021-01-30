package deadlock

import (
	"bufio"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/petermattis/goid"
)

var (
	logHeader = []byte("DEADLOCK INFO:")
)

var detector *Detector

func init() {
	detector = &Detector{
		lockers: make(map[sync.Locker]*lockerCtx),
		waiters: make(map[sync.Locker]map[int64]bool),
	}
	ticker := time.NewTicker(minDetectionTimeout / 2)
	go func() {
		for range ticker.C {
			detector.Detect()
		}
	}()
}

type lockerCtx struct {
	gid      int64
	locktime int64          // time.Now().Unix()
	waiters  map[int64]bool // key is gid
}

type Detector struct {
	mu      sync.Mutex
	lockers map[sync.Locker]*lockerCtx
	waiters map[sync.Locker]map[int64]bool // second key is gid
}

func (m *Detector) addWaiter(l sync.Locker) {
	gid := goid.Get()
	m.mu.Lock()
	if m.waiters[l] == nil {
		m.waiters[l] = make(map[int64]bool)
	}
	m.waiters[l][gid] = true
	m.mu.Unlock()
}

func (m *Detector) addLocker(l sync.Locker) {
	gid := goid.Get()
	m.mu.Lock()
	if m.waiters[l][gid] {
		delete(m.waiters[l], gid)
		m.lockers[l] = &lockerCtx{gid, time.Now().Unix(), m.waiters[l]}
	}
	m.mu.Unlock()
}

func (m *Detector) delLocker(l sync.Locker) {
	m.mu.Lock()
	delete(m.lockers, l)
	m.mu.Unlock()
}

func (m *Detector) Detect() {
	if Opts.Disable {
		return
	}
	timeout := Opts.DetectionTimeout
	if timeout < minDetectionTimeout {
		timeout = minDetectionTimeout
	}
	tiSec := int64(timeout / time.Second)

	now := time.Now().Unix()
	deads := make(map[sync.Locker]*lockerCtx)
	log.Printf("check lockers begin: %v\n", m.waiters)
	m.mu.Lock()
	for l, ctx := range m.lockers {
		if ctx.locktime+tiSec <= now {
			deads[l] = ctx
		}
	}
	log.Printf("check lockers end: %v %d %d\n", m.waiters, len(m.lockers), len(deads))
	m.mu.Unlock()
	if len(deads) > 0 {
		for l, ctx := range deads {
			m.backtrace(l, ctx)
		}
		Opts.OnDeadlock()
	}
}

func (m *Detector) backtrace(l sync.Locker, ctx *lockerCtx) {
	Opts.Logger.Write(logHeader)
	fmt.Fprintf(Opts.Logger, "the lock %p was grabbed by goroutine %d\n", l, ctx.gid)
	stacks := stacks()
	if stacks[ctx.gid] != nil {
		Opts.Logger.Write(stacks[ctx.gid])
	}
	Opts.Logger.Write([]byte("\nother waitting goroutines:\n"))
	m.mu.Lock()
	for gid := range ctx.waiters {
		if stacks[gid] != nil {
			Opts.Logger.Write(stacks[gid])
		}

	}
	m.mu.Unlock()
	if buf, ok := Opts.Logger.(*bufio.Writer); ok {
		buf.Flush()
	}
}
