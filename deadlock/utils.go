package deadlock

import (
	"bytes"
	"runtime"

	"github.com/petermattis/goid"
)

var buffer = make([]byte, 1024*16)

// Stacktraces for all goroutines.
func stacks() map[int64][]byte {
	var buf []byte
	for {
		n := runtime.Stack(buffer, true)
		if n < len(buffer) {
			buf = buffer[:n]
			break
		}
		buffer = make([]byte, 2*len(buffer))
	}
	gs := bytes.Split(buf, []byte("\n\n"))
	stacks := make(map[int64][]byte)
	for _, g := range gs {
		stacks[goid.ExtractGID(g)] = g
	}
	return stacks
}
