package cache

import (
	"testing"
)

func TestSharedGet(t *testing.T) {
	shared := newShared(1024)
	b, err := shared.get(defaultHasher.Sum64("key"), "key")
	t.Log(b, err)
}

func TestSharedSet(t *testing.T) {
	shared := newShared(1024)
	key := "key"
	value := "value"

	shared.set(defaultHasher.Sum64(key), key, string2slice(value))
	b, err := shared.get(defaultHasher.Sum64(key), key)
	t.Log(string(b), err, shared.buffer.off)

	shared.set(defaultHasher.Sum64(key), key, string2slice(value+"1"))
	b, err = shared.get(defaultHasher.Sum64(key), key)
	t.Log(string(b), err, shared.buffer.off)

	shared.set(defaultHasher.Sum64(key), key, string2slice(value+"11"))
	b, err = shared.get(defaultHasher.Sum64(key), key)
	t.Log(string(b), err, shared.buffer.off)

}

func BenchmarkSharedGet(b *testing.B) {
	var tmp []byte
	shared := newShared(1024)
	key := "key"
	value := "value"
	shared.set(defaultHasher.Sum64(key), key, string2slice(value))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = shared.get(defaultHasher.Sum64(key), key)
	}

	if tmp != nil {
		tmp[0] = 1
	}
}

func add(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i]++
		if b[i] != 0 {
			break
		}
	}
}

func BenchmarkSharedSet(b *testing.B) {
	shared := newShared(1 << 30)
	hasher := newDefaultHash()
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < b.N; i++ {
		add(k)
		err := shared.set(hasher.Sum64(slice2string(k)), slice2string(k), v)
		if err != nil {
			return
		}
	}
}

func BenchmarkSharedSetSame(b *testing.B) {
	shared := newShared(1024)
	hasher := newDefaultHash()
	key := "key"
	value := []byte("value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shared.set(hasher.Sum64(key), key, value)
	}
}
