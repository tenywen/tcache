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
	length := len(v)

	idx := s / chunkSize
	for len(v) != 0 {
		chunk := b.chunks[idx]
		if chunk == nil {
			chunk = getChunk()
		}
		v = v[copy(chunk[s%chunkSize:], v):]
		b.chunks[idx] = chunk
		idx++
	}

	if s+length > b.off {
		b.off = s + length
	}
	return nil
}

func (b *buffer) read(s, e int, dst []byte) error {
	if s > e || s < 0 || e > b.off {
		return errReadPos
	}

	endIdx := e / chunkSize
	for i := s; i <= e; {
		idx := i / chunkSize
		chunk := b.chunks[idx]
		if endIdx != idx {
			dst = append(dst, chunk[i%chunkSize:]...)
			i += (chunkSize - i%chunkSize)
		} else {
			dst = append(dst, chunk[i%chunkSize:e%chunkSize]...)
			i += (e % chunkSize)
		}
	}

	return nil
}

func (b *buffer) writeInt(s int, v int) error {
	small := getSmall()
	defer recycleSmall(small)
	length := encode(int(v), small)
	return b.write(s, small[:length])
}
