package cache

import (
	"encoding/binary"
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
	if n+b.off <= cap(b.bytes) {
		b.bytes = b.bytes[:b.off+n]
		return true
	}

	return false
}

func (b *buffer) rewrite(s, n int, k, v []byte) {
	if s+n > b.off {
		panic("buffer rewrite fatal")
	}

	s += writeInt(b.bytes[s:], n, len(k), len(v))
	write(b.bytes[s:], k, v)
}

func (b *buffer) write(ps ...[]byte) {
	var l int
	for _, p := range ps {
		l += len(p)
	}

	if !b.grow(l) {
		return
	}

	b.off += write(b.bytes[b.off:], ps...)
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

func (b buffer) btoi(si int) int {
	if si > b.off || si+int64Len > b.off {
		return undefined
	}

	data, err := b.read(si, si+int64Len)
	if err != nil {
		return undefined
	}

	return int(binary.LittleEndian.Uint64(data))
}

func (b *buffer) writeInt(vs ...int) {
	if !b.grow(int64Len * len(vs)) {
		return
	}
	b.off += writeInt(b.bytes[b.off:], vs...)
}

func writeInt(bytes []byte, vs ...int) int {
	var start int
	for _, v := range vs {
		binary.LittleEndian.PutUint64(bytes[start:], uint64(v))
		start += int64Len
	}
	return start
}

func write(bytes []byte, ps ...[]byte) int {
	var start int
	for _, p := range ps {
		start += copy(bytes[start:], p)
	}
	return start
}
