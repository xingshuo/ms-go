package deadlock

import (
	"sync"
	"time"
)

var (
	logHeader = []byte("DEADLOCK INFO:")
)

var detector *Detector

func init() {
	detector = &Detector{
		lockers:   make(map[*Mutex]bool),
		rwlockers: make(map[*RWMutex]bool),
	}
	ticker := time.NewTicker(minDetectionTimeout / 2)
	go func() {
		for range ticker.C {
			detector.Detect()
		}
	}()
}

type lockCtx struct {
	gid      int64
	locktime int64
}

type Detector struct {
	mu        sync.Mutex
	lockers   map[*Mutex]bool
	rwmu      sync.Mutex
	rwlockers map[*RWMutex]bool
}

func (d *Detector) addLocker(l *Mutex) {
	d.mu.Lock()
	d.lockers[l] = true
	d.mu.Unlock()
}

func (d *Detector) delLocker(l *Mutex) {
	d.mu.Lock()
	delete(d.lockers, l)
	d.mu.Unlock()
}

func (d *Detector) addRWLocker(l *RWMutex) {
	d.rwmu.Lock()
	d.rwlockers[l] = true
	d.rwmu.Unlock()
}

func (d *Detector) delRWLocker(l *RWMutex) {
	d.rwmu.Lock()
	delete(d.rwlockers, l)
	d.rwmu.Unlock()
}

func (d *Detector) Detect() {
	if Opts.Disable {
		return
	}
	timeout := Opts.DetectionTimeout
	if timeout < minDetectionTimeout {
		timeout = minDetectionTimeout
	}
	tiSec := int64(timeout / time.Second)
	now := time.Now().Unix()
	stacks := stacks()

	d.mu.Lock()
	for l := range d.lockers {
		go l.detect(now, tiSec, stacks)
	}
	d.mu.Unlock()

	d.rwmu.Lock()
	for l := range d.rwlockers {
		go l.detect(now, tiSec, stacks)
	}
	d.rwmu.Unlock()
}
