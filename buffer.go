package cache

import (
	"errors"
)

const (
	int64Len = 8
)

var (
	errReadPos  = errors.New("read start pos or end pos error")
	errCapLimit = errors.New("buffer cap exceeds limit ")
)

type buffer struct {
	off   int
	bytes []byte
	pools pools
}

func newBuffer(min, max int) buffer {
	return buffer{
		pools: newPools(min, max),
	}
}

func (b *buffer) grow(n int) bool {
	if n <= 0 {
		return true
	}

	if b.tryGrowByReslice(n) {
		return true
	}

	c := cap(b.bytes) + n
	p2 := power2(c)
	if p2 < c {
		return false
	}

	m := b.pools.get(p2)
	copy(m, b.bytes[:b.off])
	b.pools.put(b.bytes)
	b.bytes = m
	return true
}

func (b *buffer) tryGrowByReslice(n int) bool {
	if len(b.bytes) >= n+b.off {
		return true
	}

	if n+b.off <= cap(b.bytes) {
		b.bytes = b.bytes[:b.off+n]
		return true
	}

	return false
}

func (b *buffer) write(si, n int, k, v []byte) {
	if si+n > b.off {
		b.grow(si + n - b.off)
	}

	copy(b.bytes[si:], k)
	copy(b.bytes[si+len(k):], v)
	if si+n > b.off {
		b.off = si + n
	}
}

func (b buffer) read(start, end int) ([]byte, error) {
	if start < 0 || end < start {
		return nil, errReadPos
	}

	if end > b.off {
		return nil, errCapLimit
	}

	return b.bytes[start:end], nil
}
