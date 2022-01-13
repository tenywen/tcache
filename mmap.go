// +build !appengine,!windows

package cache

import (
	"fmt"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	chunkSize      = 1 << 16
	chunksPerAlloc = 512
)

var (
	chunks []*[chunkSize]byte
	mu     sync.Mutex
)

func mmap() {
	// Allocate offheap memory, so GOGC won't take into account cache size.
	// This should reduce free memory waste.
	data, err := unix.Mmap(-1, 0, chunkSize*chunksPerAlloc, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANONYMOUS|unix.MAP_PRIVATE)
	if err != nil {
		panic(fmt.Errorf("cannot allocate %d bytes via mmap: %s", chunkSize*chunksPerAlloc, err))
	}
	for len(data) > 0 {
		array := (*[chunkSize]byte)(unsafe.Pointer(&data[0]))
		chunks = append(chunks, array)
		data = data[chunkSize:]
	}
}

func getChunk() []byte {
	//return make([]byte, chunkSize)
	mu.Lock()
	if len(chunks) == 0 {
		mmap()
	}
	n := len(chunks) - 1
	p := chunks[n]
	chunks[n] = nil
	chunks = chunks[:n]
	mu.Unlock()
	return p[:]
}

func putChunk(chunk []byte) {
	if chunk == nil {
		return
	}

	chunk = chunk[:chunkSize]
	p := (*[chunkSize]byte)(unsafe.Pointer(&chunk[0]))

	mu.Lock()
	chunks = append(chunks, p)
	mu.Unlock()
}
