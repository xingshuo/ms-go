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
	Logger           io.Writer
}{
	Disable:          false,
	DetectionTimeout: minDetectionTimeout,
	OnDeadlock:       defaultOnDeadlock,
	Logger:           os.Stderr,
}
