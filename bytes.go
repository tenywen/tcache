package cache

import (
	"encoding/binary"
	"errors"
)

const (
	int64Len = 8
)

var (
	errReadPos  = errors.New("read start_pos or end_pos error")
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

	c := int64(cap(b.m)) + n
	p2 := power2(c)
	if p2 < c {
		return false
	}

	m := make([]byte, p2)
	copy(m, b.m[:b.off])
	b.m = m
	return true
}

func (b *buffer) tryGrowByReslice(n int64) bool {
	if n+b.off <= int64(cap(b.m)) {
		b.m = b.m[:b.off+n]
		return true
	}

	return false
}

func (b *buffer) write(p []byte) int64 {
	l := int64(len(p))
	if !b.grow(l) {
		return undefined
	}
	copy(b.m[b.off:], p)
	b.off += l
	return l
}

func (b buffer) read(start, end int64) ([]byte, error) {
	if start < 0 || end < start {
		return nil, errReadPos
	}

	if end > b.off {
		println("end:", end, "off:", b.off)
		return nil, errCapLimit
	}

	data := make([]byte, end-start+1)
	copy(data, b.m[start:end+1])
	return data, nil
}

func (b buffer) int64(start, end int64) int64 {
	data, err := b.read(start, end)
	if err != nil {
		return undefined
	}

	return int64(binary.LittleEndian.Uint64(data))
}

func (b buffer) tail() int64 {
	return b.off
}

func (b *buffer) putInt64(v int64) {
	if !b.grow(int64Len) {
		return
	}
	binary.LittleEndian.PutUint64(b.m[b.off:], uint64(v))
	b.off += int64Len
}
