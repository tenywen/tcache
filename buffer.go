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

	m := make([]byte, p2)
	copy(m, b.bytes[:b.off])
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
	if s+n > b.off || len(k)+len(v) != n {
		panic("buffer rewrite fatal")
	}

	if len(k)+len(v) != n {
		panic("buffer rewrite size fatal")
	}

	write(b.bytes[s:], k)
	write(b.bytes[s+len(k):], v)
}

func (b *buffer) write(p []byte) {
	l := len(p)
	if !b.grow(l) {
		return
	}

	b.off += write(b.bytes[b.off:], p)
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

func (b *buffer) writeInt(v int) {
	if !b.grow(int64Len) {
		return
	}
	b.off += writeInt(b.bytes[b.off:], v)
}

func writeInt(bytes []byte, v int) int {
	binary.LittleEndian.PutUint64(bytes, uint64(v))
	return int64Len
}

func write(bytes []byte, p []byte) int {
	return copy(bytes, p)
}
