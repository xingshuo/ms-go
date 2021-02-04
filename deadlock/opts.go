package deadlock

import (
	"io"
	"os"
	"time"
)

const (
	minDetectionTimeout = 10 * time.Second
)

var (
	defaultOnDeadlock = func() {
		os.Exit(2)
	}
)

var Opts = struct {
	Disable          bool
	DetectionTimeout time.Duration
	OnDeadlock       func()
	PrintDeadlockRoutineNum uint32 // 发生死锁时, 打印相关goroutine数量, 默认打印全部
	Logger           io.Writer
}{
	Disable:          false,
	DetectionTimeout: minDetectionTimeout,
	OnDeadlock:       defaultOnDeadlock,
	PrintDeadlockRoutineNum: 0,
	Logger:           os.Stderr,
}

func isValidPrintStack(gNum uint32) bool {
	if Opts.PrintDeadlockRoutineNum == 0 {
		return true
	}
	return gNum <=  Opts.PrintDeadlockRoutineNum
}
