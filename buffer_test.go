package cache

import (
	"testing"
)

func TestGrow(t *testing.T) {
	b := newBuffer(1024)
	b.grow(2)
	t.Log(len(b.bytes), cap(b.bytes))
	b.grow(23)
	t.Log(len(b.bytes), cap(b.bytes))
	b.grow(125)
	t.Log(len(b.bytes), cap(b.bytes))
	b.grow(250)
	t.Log(len(b.bytes), cap(b.bytes))
	b.grow(257)
	t.Log(len(b.bytes), cap(b.bytes))
}

func TestWriteBytes(t *testing.T) {
	b := newBuffer(1024)
	v := []byte("2222")
	k := []byte("111")
	b.write(b.off, len(v)+len(v), v, v)
	t.Log(b.off)
	b.write(0, len(k)+len(v), k, k)
	t.Log(b.off)
}

func TestReadBytes(t *testing.T) {
	b := newBuffer(1024)
	k := []byte("key234234")
	v := []byte("value34234")
	b.write(b.off, len(k)+len(v), k, v)
	b.write(b.off, len(k)+len(v), k, v)
	data, err := b.read(b.off-len(v), b.off)
	t.Log(string(data), err)
	data, err = b.read(b.off-len(v)-len(k), b.off-len(v))
	t.Log(string(data), err)
}
