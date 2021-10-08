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
	off   int64
	bytes []byte
}

func (b *buffer) grow(n int64) bool {
	if n <= 0 {
		return true
	}

	if b.tryGrowByReslice(n) {
		return true
	}

	c := int64(cap(b.bytes)) + n
	p2 := power2(c)
	if p2 < c {
		return false
	}

	m := make([]byte, p2)
	copy(m, b.bytes[:b.off])
	b.bytes = m
	return true
}

func (b *buffer) tryGrowByReslice(n int64) bool {
	if n+b.off <= int64(cap(b.bytes)) {
		b.bytes = b.bytes[:b.off+n]
		return true
	}

	return false
}

func (b *buffer) write(p []byte) int64 {
	l := int64(len(p))
	if !b.grow(l) {
		return undefined
	}
	copy(b.bytes[b.off:], p)
	b.off += l
	return l
}

/*
func (b buffer) copyRead(start, end int64) ([]byte, error) {
	if src, err := b.read(start, end); err != nil {
		return nil, err
	}

	return nil, nil
}
*/

func (b buffer) read(start, end int64) ([]byte, error) {
	if start < 0 || end < start {
		return nil, errReadPos
	}

	if end > b.off {
		return nil, errCapLimit
	}

	return b.bytes[start:end], nil
}

func (b buffer) btoi(si int64) int64 {
	if si > b.off || si+int64Len > b.off {
		return undefined
	}

	data, err := b.read(si, si+int64Len)
	if err != nil {
		return undefined
	}

	return int64(binary.LittleEndian.Uint64(data))

}

func (b *buffer) writeInt64(v int64) {
	if !b.grow(int64Len) {
		return
	}
	binary.LittleEndian.PutUint64(b.bytes[b.off:], uint64(v))
	b.off += int64Len
}
