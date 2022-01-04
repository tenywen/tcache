package cache

import (
	"testing"
)

func TestBufferWrite(t *testing.T) {
	buffer := newBuffer()
	tmp := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	buffer.write(buffer.off, tmp)
	t.Logf("%s %d %d", buffer.chunks, buffer.off, len(buffer.chunks))
	tmp = []byte{'a', 'b', 'c', 'a', 'b', 'c', 'a', 'b', 'c', 'a', 'b', 'c'}
	buffer.write(buffer.off, tmp)
	t.Logf("%v %d %d", buffer.chunks, buffer.off, len(buffer.chunks))
}

func TestBufferRead(t *testing.T) {
	buffer := newBuffer()

	tmp := []byte{'a', 'b', 'c', 'a', 'b', 'c', 'a', 'b', 'c', 'a', 'b', 'c'}
	buffer.write(buffer.off, tmp)

	s := 7
	e := buffer.off - 1
	var dst []byte
	dst, err := buffer.read(s, e, dst)
	t.Log(err)
	t.Logf("%s\n", dst)
}

func BenchmarkBufferWrite(b *testing.B) {
	const items = 1 << 16
	buffer := newBuffer()
	b.ReportAllocs()
	b.SetBytes(items)
	var s int
	for i := 0; i < b.N; i++ {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for i := 0; i < items; i++ {
			k[0]++
			if k[0] == 0 {
				k[1]++
			}
			buffer.write(s, k)
			s += len(k)
			buffer.write(s, v)
			s += len(v)
		}
	}

	b.Logf("%v %d \n", buffer.chunks[0][:20], len(buffer.chunks))
}
