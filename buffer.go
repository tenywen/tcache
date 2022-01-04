package cache

import (
	"errors"
)

const (
	int64Len = 8
)

var (
	errReadPos    = errors.New("read start pos or end pos error")
	errCapLimit   = errors.New("buffer cap exceeds limit ")
	errNotContent = errors.New("not content in this start pos")
)

type buffer struct {
	off    int
	chunks [][]byte
}

func newBuffer() buffer {
	return buffer{}
}

func (b *buffer) write(s int, v []byte) error {
	if s > b.off {
		return errReadPos
	}

	length := len(v)
	offset := s % chunkSize
	idx := s / chunkSize
	for len(v) != 0 {
		if len(b.chunks) == idx {
			b.chunks = append(b.chunks, getChunk())
		}

		chunk := b.chunks[idx]
		v = v[copy(chunk[offset:], v):]
		b.chunks[idx] = chunk
		offset = 0
		idx++
	}

	if s+length > b.off {
		b.off = s + length
	}
	return nil
}

func (b *buffer) read(s, e int, dst []byte) ([]byte, error) {
	if s > e || s < 0 || e > b.off {
		return nil, errReadPos
	}

	length := e - s

	if dst == nil {
		dst = make([]byte, length)
	}

	if len(dst) > length {
		dst = dst[:length]
	}

	tmp := dst
	offset := s % chunkSize
	idx := s / chunkSize

	for len(tmp) > 0 {
		chunk := b.chunks[idx]
		tmp = tmp[copy(tmp, chunk[offset:]):]
		offset = 0
		idx++
	}

	return dst, nil
}

func (b *buffer) writeInt(s int, v int) error {
	small := getSmall()
	defer recycleSmall(small)
	length := encode(int(v), small)
	return b.write(s, small[:length])
}
