package cache

import "testing"

func TestNewPools(t *testing.T) {
	defer func() {
		recover()
	}()
	newPools(1, 1)
	newPools(1, 16)
	newPools(1, 256)
	newPools(1, 1<<16)
	newPools(1<<16, 1)
}

func TestGetPool(t *testing.T) {
	ps := newPools(8, 1<<16)
	ps.getPool(1)
	ps.getPool(2)
	ps.getPool(4)
	ps.getPool(8)
	ps.getPool(16)
	ps.getPool(32)
	ps.getPool(64)
	ps.getPool(128)
	ps.getPool(256)
	ps.getPool(512)
	ps.getPool(1024)
	ps.getPool(2043)
	ps.getPool(2048)
	ps.getPool(2049)
}

func TestGet(t *testing.T) {
	ps := newPools(1, 1<<20)
	ps.get(1 << 20)
}

func BenchmarkGetWithoutPut(b *testing.B) {
	var bytes []byte
	ps := newPools(1, b.N)
	for i := 1; i <= b.N; i++ {
		bytes = ps.get(i)
	}

	bytes[0] = byte(1)
}

func BenchmarkGet(b *testing.B) {
	var bytes []byte
	ps := newPools(1, b.N)
	for i := 1; i <= b.N; i++ {
		bytes = ps.get(i)
		bytes[0] = byte(1)
		ps.put(bytes)
	}

	bytes[0] = byte(1)
}
