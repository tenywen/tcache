package cache

import (
	"strconv"
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

func BenchmarkSharedSet(b *testing.B) {
	shared := newShared(1 << 30)
	hasher := newDefaultHash()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.FormatInt(int64(i), 10)
		err := shared.set(hasher.Sum64(key), key, string2slice(key))
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
